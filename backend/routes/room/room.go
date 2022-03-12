package room

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/taiprogramer/simple-poker-game/backend/db"
	"github.com/taiprogramer/simple-poker-game/backend/routes"
	"github.com/taiprogramer/simple-poker-game/backend/secure"
)

type UserInRoomSchema struct {
	ID    uint `json:"id"`
	Ready bool `json:"ready"`
}

type RoomSchemaResponse struct {
	ID      uint               `json:"id"`
	Code    string             `json:"code"`
	Playing bool               `json:"playing"`
	Private bool               `json:"private"`
	Owner   uint               `json:"owner"`
	Table   uint               `json:"table"`
	Users   []UserInRoomSchema `json:"users"`
}

type RoomBody struct {
	UserID   uint   `json:"user_id"`
	Password string `json:"password"`
}

func isAuthenticatedUser(userID uint, username interface{}) bool {
	var user db.User
	db.DB.First(&user, userID)
	return user.Username == username
}

// TODO
func getNextUniqueCode() string {
	return "[must be unique for each room]"
}

func createNewRoom(userID uint, private bool, password string) (*RoomSchemaResponse, bool) {
	code := getNextUniqueCode()
	room := db.Room{
		Code:     code,
		Playing:  false,
		Private:  private,
		Password: password,
		UserID:   userID,
	}
	result := db.DB.Create(&room)
	if result.RowsAffected == 0 {
		return &RoomSchemaResponse{}, false
	}
	roomResponse := RoomSchemaResponse{
		ID:      room.ID,
		Code:    room.Code,
		Playing: room.Playing,
		Private: room.Private,
		Owner:   room.UserID,
		Users:   []UserInRoomSchema{},
	}
	return &roomResponse, true
}

// CreateNewRoomHandler use for creating new room with or without password
func CreateNewRoomHandler(c *fiber.Ctx) error {
	body := new(RoomBody)

	if err := c.BodyParser(body); err != nil {
		e := routes.NewErrorResponse([]string{"Unknown error"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	// invalid user
	if body.UserID == 0 {
		e := routes.NewErrorResponse([]string{"User id is not valid."})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	// check if user_id is id of authenticated user.
	usr := secure.GetJWTClaim("usr", c.Locals("token"))
	if !isAuthenticatedUser(body.UserID, usr) {
		e := routes.NewErrorResponse([]string{"Are you hacker?"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	private := body.Password != ""
	response, ok := createNewRoom(body.UserID, private, body.Password)

	if !ok {
		e := routes.NewErrorResponse([]string{"Can not create new room."})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func getRoomById(id int) (*RoomSchemaResponse, bool) {
	var room db.Room
	result := db.DB.First(&room, id)
	if result.RowsAffected == 0 {
		return &RoomSchemaResponse{}, false
	}
	// find table belongs to room
	var table db.Table
	db.DB.Where("room_id = ?", id).First(&table)

	// users in room
	usersInRoom := make([]UserInRoomSchema, 0)
	var waitingLists []db.WaitingList
	db.DB.Where("room_id", id).Find(&waitingLists)
	for _, v := range waitingLists {
		usersInRoom = append(usersInRoom, UserInRoomSchema{
			ID:    v.UserID,
			Ready: v.Ready,
		})
	}

	roomResponse := RoomSchemaResponse{
		ID:      room.ID,
		Code:    room.Code,
		Playing: room.Playing,
		Private: room.Private,
		Users:   usersInRoom,
		Owner:   room.UserID,
		Table:   table.ID,
	}

	return &roomResponse, true
}

func GetSpecificRoomHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		e := routes.NewErrorResponse([]string{"Please supply a room id!"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	response, ok := getRoomById(id)

	if !ok {
		e := routes.NewErrorResponse([]string{"Room not found."})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
