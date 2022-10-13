package socket_mgmt

import "github.com/gofiber/websocket/v2"

type SocketConnection struct {
	userID int
	c      *websocket.Conn
}

var roomSocketConnections map[int][]SocketConnection = make(map[int][]SocketConnection)

func BroadcastMsgToRoom(msg string, roomID int) {
	connections, ok := roomSocketConnections[roomID]
	if ok {
		for _, conn := range connections {
			conn.c.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}

func UnicastMsgToUser(msg string, roomID int, userID int) {
	connections, ok := roomSocketConnections[roomID]
	if ok {
		for _, conn := range connections {
			if conn.userID == userID {
				conn.c.WriteMessage(websocket.TextMessage, []byte(msg))
				break
			}
		}
	}
}

func StoreSocketConnection(c *websocket.Conn, userID int, roomID int) {
	socketConn := SocketConnection{
		c:      c,
		userID: userID,
	}

	roomSocketConnections[roomID] = append(roomSocketConnections[roomID], socketConn)
}

func RemoveSocketConnection(userID int, roomID int) {
	i := 0
	for _, socketConnection := range roomSocketConnections[roomID] {
		if socketConnection.userID != userID {
			roomSocketConnections[roomID][i] = socketConnection
			i++
		}
	}
	roomSocketConnections[roomID] = roomSocketConnections[roomID][:i]
}
