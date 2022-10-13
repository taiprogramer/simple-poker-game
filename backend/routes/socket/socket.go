package socket

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/websocket/v2"
	"github.com/taiprogramer/simple-poker-game/backend/db"
	roomRepo "github.com/taiprogramer/simple-poker-game/backend/repo/room"
	tableRepo "github.com/taiprogramer/simple-poker-game/backend/repo/table"
	userTableCardRepo "github.com/taiprogramer/simple-poker-game/backend/repo/user_table_card"
	"github.com/taiprogramer/simple-poker-game/backend/routes/socket/room_card"
	"github.com/taiprogramer/simple-poker-game/backend/routes/socket/socket_mgmt"
)

// map userID to roomID
var userInRooms map[int]int = make(map[int]int)

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

	room_card.ShuffleCards(roomID)
	for _, v := range waitings {
		// add players to tables.
		if v.Ready {
			tableID := tableRepo.CreateNewTable(v.UserID, room.ID, 0, 0, false)
			// deal 2 cards for playing users.
			card1 := room_card.DealNextCard(roomID)
			card2 := room_card.DealNextCard(roomID)
			card1ID := room_card.GetCardID(card1)
			card2ID := room_card.GetCardID(card2)
			userTableCardRepo.AddNewCard(tableID, v.UserID, card1ID)
			userTableCardRepo.AddNewCard(tableID, v.UserID, card2ID)
			socket_mgmt.UnicastMsgToUser("table="+fmt.Sprint(tableID), roomID, int(v.UserID))
		}
	}
	room.Playing = true
	roomRepo.UpdateRoom(room)
	socket_mgmt.BroadcastMsgToRoom("the game is started", roomID)
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
