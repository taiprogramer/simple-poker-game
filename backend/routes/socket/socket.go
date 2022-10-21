package socket

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/websocket/v2"
	"github.com/taiprogramer/simple-poker-game/backend/db"
	betHistoryRepo "github.com/taiprogramer/simple-poker-game/backend/repo/bet_histories"
	roomRepo "github.com/taiprogramer/simple-poker-game/backend/repo/room"
	tableRepo "github.com/taiprogramer/simple-poker-game/backend/repo/table"
	userTableCardRepo "github.com/taiprogramer/simple-poker-game/backend/repo/user_table_card"
	"github.com/taiprogramer/simple-poker-game/backend/routes/socket/room_card"
	"github.com/taiprogramer/simple-poker-game/backend/routes/socket/socket_mgmt"
)

// map userID to roomID
var userInRooms map[int]int = make(map[int]int)

type UserTurn struct {
	UserID           int
	HasPerformAction bool

	// total amount money of the current round, the round is only finished if
	// all users have the same amount and they've already performed their
	// action.
	Amount int
}

// track current turn
var tableUsersTurn map[int][]UserTurn = make(map[int][]UserTurn)

func startNewGame(room *db.Room, userID int) {
	roomID := int(room.ID)
	// not the owner can not issue "start" command
	if room.UserID != uint(userID) {
		return
	}
	/* Set up Stage 0
	- create new table.
	- add ready users to table.
	- deal 2 cards for playing users.
	*/
	waitings := roomRepo.GetWaitingListsByRoomID(roomID)

	tableID := tableRepo.CreateNewTable(uint(userID), room.ID, 1, 0, false)
	room_card.ShuffleCards(roomID)
	for _, v := range waitings {
		// add players to tables.
		if v.Ready {
			// deal 2 cards for playing users.
			card1 := room_card.DealNextCard(roomID)
			card2 := room_card.DealNextCard(roomID)
			card1ID := room_card.GetCardID(card1)
			card2ID := room_card.GetCardID(card2)
			userTableCardRepo.AddNewCard(tableID, v.UserID, card1ID)
			userTableCardRepo.AddNewCard(tableID, v.UserID, card2ID)
			// track current turn
			userTurn := UserTurn{
				UserID:           int(v.UserID),
				HasPerformAction: false,
				Amount:           0,
			}
			tableUsersTurn[int(tableID)] = append(tableUsersTurn[int(tableID)], userTurn)
			socket_mgmt.UnicastMsgToUser("table="+fmt.Sprint(tableID), roomID, int(v.UserID))
		}
	}
	room.Playing = true
	// big blind and small blind perform action
	for i := 0; i < 2; i++ {
		userTurn := tableUsersTurn[int(tableID)][i]
		userTurn.HasPerformAction = true
		userTurn.Amount = 100 * (i + 1)
		tableUsersTurn[int(tableID)][i] = userTurn
		actionID := betHistoryRepo.FindActionIDByName("raise")
		betHistoryRepo.WriteBetHistory(int(tableID), userTurn.UserID, actionID, 100*(i+1), 1)
	}
	// next current is third user
	table := tableRepo.FindTableByRoomID(room.ID)
	table.Pot = 300
	table.UserID = uint(tableUsersTurn[int(tableID)][2].UserID)
	tableRepo.UpdateTable(&table)

	roomRepo.UpdateRoom(room)
	socket_mgmt.BroadcastMsgToRoom("the game is started", roomID)
}

func performActionPostHandler(userID, roomID int) {
	isNextRound := true
	table := tableRepo.FindTableByRoomID(uint(roomID))
	usersTurns := tableUsersTurn[int(table.ID)]
	betHistory := betHistoryRepo.GetLatest(userID)
	table.Pot += betHistory.Amount
	totalAmount := betHistoryRepo.GetTotalAmountByRoundAndUserID(int(table.ID), table.Round, userID)

	for i, v := range usersTurns {
		if v.Amount < totalAmount {
			v.HasPerformAction = false
		}
		if v.UserID == userID {
			v.HasPerformAction = true
			v.Amount = totalAmount
		}
		usersTurns[i] = v
	}

	// find next current turn
	nextCurrentTurnIndex := 0
	for i, v := range usersTurns {
		if !v.HasPerformAction {
			isNextRound = false
			nextCurrentTurnIndex = i
			break
		}
	}
	// re-assign the current turn
	table.UserID = uint(usersTurns[nextCurrentTurnIndex].UserID)

	if isNextRound {
		table.Round += 1
		// reset users current turn
		for i := 0; i < len(usersTurns); i++ {
			usersTurns[i].HasPerformAction = false
			usersTurns[i].Amount = 0
		}
	}
	tableRepo.UpdateTable(&table)
	socket_mgmt.BroadcastMsgToRoom("table="+fmt.Sprint(table.ID), roomID)
}

func SocketHandler(c *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)

	userID, _ := strconv.Atoi(c.Params("user_id"))
	roomID, _ := strconv.Atoi(c.Query("room"))
	userInRooms[userID] = roomID

	// get the room
	room, _ := roomRepo.FindRoomByID(roomID)

	// notify for existing players when new user join room
	socket_mgmt.BroadcastMsgToRoom("new user join room", roomID)

	socket_mgmt.StoreSocketConnection(c, userID, roomID)

	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			break
		}
		command := string(msg[:])
		// start the game, after start, the game is in Stage 0
		if strings.Compare(command, "start") == 0 {
			startNewGame(room, userID)
		}
		if strings.Compare(command, "ready") == 0 {
			socket_mgmt.BroadcastMsgToRoom("room status was changed", roomID)
		}
		if strings.Compare(command, "has performed action") == 0 {
			performActionPostHandler(userID, roomID)
		}

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}

	socket_mgmt.RemoveSocketConnection(userID, roomID)

	// when users leave the room (connection is closed), if they are the
	// owner of the room, delete the room and transfer that room to the new
	// owner.
	roomID, roomExists := userInRooms[userID]
	if !roomExists {
		return
	}

	if room.Playing {
		// TODO: add logic when users leave the room while playing status is
		// being true.
		return
	}
	// user is owner and room playing status is not true.
	// find waiting lists
	waiting, _ := roomRepo.FindWaitingListsByRoomID(roomID)
	if room.UserID == uint(userID) {
		if len(*waiting) == 1 {
			// safe to delete the entire room
			roomRepo.DeleteRoom(room)
			roomRepo.DeleteWaitingListsByRoomID(roomID)
			return
		}
		// assign room to random remain users
		for _, v := range *waiting {
			if v.UserID != uint(userID) {
				room.UserID = v.UserID
				roomRepo.UpdateRoom(room)
				break
			}
		}
	}
	// remove user from waiting list
	for _, v := range *waiting {
		if v.UserID == uint(userID) {
			roomRepo.DeleteWaitingListByUserID(userID, &v)
			break
		}
	}

	// notify room status was changed when players leave room
	socket_mgmt.BroadcastMsgToRoom("room status was changed", roomID)
}
