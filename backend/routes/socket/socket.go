package socket

import (
	"log"
	"strconv"

	"github.com/gofiber/websocket/v2"
	roomRepo "github.com/taiprogramer/simple-poker-game/backend/repo/room"
)

// map userID to roomID
var userInRooms map[int]int = make(map[int]int)

func SocketHandler(c *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)

	userID, _ := strconv.Atoi(c.Params("user_id"))
	roomID, _ := strconv.Atoi(c.Query("room"))
	userInRooms[userID] = roomID

	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
	// when users leave the room (connection is closed), if they are the
	// owner of the room, delete the room and transfer that room to the new
	// owner.
	roomID, roomExists := userInRooms[userID]
	if !roomExists {
		return
	}

	// get the room
	room, _ := roomRepo.FindRoomByID(roomID)
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
