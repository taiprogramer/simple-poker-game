package socket

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/websocket/v2"
	cardRepo "github.com/taiprogramer/simple-poker-game/backend/repo/card"
	roomRepo "github.com/taiprogramer/simple-poker-game/backend/repo/room"
	tableRepo "github.com/taiprogramer/simple-poker-game/backend/repo/table"
	userTableCardRepo "github.com/taiprogramer/simple-poker-game/backend/repo/user_table_card"
)

// map userID to roomID
var userInRooms map[int]int = make(map[int]int)

// each room will have available cards
// cards encoded:
// first character is number
// second character is suit
// 13 -> A of SPADE
// 22 -> 2 of CLUB
// 30 -> 3 of DIAMOND
// 41 -> 4 of HEART
// ..
// special cases:
// T1 -> 10 of HEART
// J1 -> J of HEART
var roomCards map[int][]string = make(map[int][]string)

type SocketConnection struct {
	userID int
	c      *websocket.Conn
}

var roomSocketConnections map[int][]SocketConnection = make(map[int][]SocketConnection)

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func shuffleCards(roomID int) {
	var cards []string = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
	var suits []string = []string{"0", "1", "2", "3"}
	var deckOfCards []string = make([]string, 0)
	// generate new deck of cards
	for _, card := range cards {
		for _, suit := range suits {
			finalCard := card + suit
			deckOfCards = append(deckOfCards, finalCard)
		}
	}
	// shuffle deck of cards
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 52; i++ {
		index := randInt(0, 52)
		temp := deckOfCards[i]
		deckOfCards[i] = deckOfCards[index]
		deckOfCards[index] = temp
	}

	roomCards[roomID] = deckOfCards
}

func dealNextCard(roomID int) string {
	deckOfCards := roomCards[roomID]
	card := deckOfCards[len(deckOfCards)-1]
	roomCards[roomID] = deckOfCards[:len(deckOfCards)-1]
	return card
}

func getCardID(card string) uint {
	cardNumber := map[string]int{
		"T": 10,
		"J": 11,
		"Q": 12,
		"K": 13,
	}
	number, ok := cardNumber[string(card[0])]
	if !ok {
		number, _ = strconv.Atoi(string(card[0]))
	}
	suit, _ := strconv.Atoi(string(card[1]))
	return cardRepo.FindCardIDByNumberAndSuit(number, suit)
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

	connections, ok := roomSocketConnections[roomID]
	if ok {
		// notify for existing players when new user join room
		for _, conn := range connections {
			conn.c.WriteMessage(websocket.TextMessage, []byte("new user join room"))
		}
	}

	socketConn := SocketConnection{
		c:      c,
		userID: userID,
	}

	roomSocketConnections[roomID] = append(roomSocketConnections[roomID], socketConn)

	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			break
		}
		command := string(msg[:])
		// start the game, after start, the game is in Stage 0
		if strings.Compare(command, "start") == 0 {
			// not the owner can not issue "start" command
			if room.UserID != uint(userID) {
				continue
			}
			/* Set up Stage 0
			- create new table.
			- add ready users to table.
			- deal 2 cards for playing users.
			*/
			waitings := roomRepo.GetWaitingListsByRoomID(roomID)

			shuffleCards(roomID)
			for _, v := range waitings {
				// add players to tables.
				if v.Ready {
					tableID := tableRepo.CreateNewTable(v.UserID, room.ID, 0, 0, false)
					// deal 2 cards for playing users.
					card1 := dealNextCard(roomID)
					card2 := dealNextCard(roomID)
					card1ID := getCardID(card1)
					card2ID := getCardID(card2)
					userTableCardRepo.AddNewCard(tableID, uint(userID), card1ID)
					userTableCardRepo.AddNewCard(tableID, uint(userID), card2ID)
				}
			}
		}

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}

	// remove socket connection
	i := 0
	for _, socketConnection := range roomSocketConnections[roomID] {
		if socketConnection.userID != userID {
			roomSocketConnections[roomID][i] = socketConnection
			i++
		}
	}
	roomSocketConnections[roomID] = roomSocketConnections[roomID][:i]

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
}
