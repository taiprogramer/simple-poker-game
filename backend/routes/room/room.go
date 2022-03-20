package room

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/taiprogramer/simple-poker-game/backend/db"
	"github.com/taiprogramer/simple-poker-game/backend/routes"
	"github.com/taiprogramer/simple-poker-game/backend/secure"
	"gorm.io/gorm"
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

	// add room owner to waiting list
	waitingList := db.WaitingList{
		UserID:         userID,
		AvailableMoney: 0,
		Ready:          true,
		RoomID:         room.ID,
	}

	db.DB.Create(&waitingList)

	roomResponse := RoomSchemaResponse{
		ID:      room.ID,
		Code:    room.Code,
		Playing: room.Playing,
		Private: room.Private,
		Owner:   room.UserID,
		Users: []UserInRoomSchema{
			{
				ID:    waitingList.UserID,
				Ready: waitingList.Ready,
			},
		},
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

func getListRooms(offset int, limit int) ([]RoomSchemaResponse, bool) {
	var rooms []db.Room
	db.DB.Offset(offset).Limit(limit).Find(&rooms)

	listRooms := make([]RoomSchemaResponse, 0)
	for _, room := range rooms {
		var table db.Table
		db.DB.Where("room_id = ?", room.ID).First(&table)
		usersInRoom := make([]UserInRoomSchema, 0)
		var waitingLists []db.WaitingList
		db.DB.Where("room_id", room.ID).Find(&waitingLists)
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
		listRooms = append(listRooms, roomResponse)
	}

	return listRooms, true
}

func GetListRoomsHandler(c *fiber.Ctx) error {
	offset, offsetErr := strconv.Atoi(c.Query("offset"))
	limit, limitErr := strconv.Atoi(c.Query("limit"))
	if offsetErr != nil || limitErr != nil {
		e := routes.NewErrorResponse([]string{"Missing query: limit " +
			"and offset are required"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	if offset < 0 || limit < 1 || limit > 100 {
		e := routes.NewErrorResponse([]string{
			"Invalid query: limit in [1..100] and offset >= 0",
		})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	listRooms, _ := getListRooms(offset, limit)
	return c.Status(fiber.StatusOK).JSON(listRooms)
}

func deleteRoomByID(id int, userID uint) error {
	var room db.Room
	db.DB.Where("id = ?", id).First(&room)
	if room.ID == 0 {
		return errors.New("Room not found")
	}
	if room.UserID != userID {
		return errors.New("You are not the room owner")
	}
	if room.Playing {
		return errors.New("Room is in playing state")
	}

	result := db.DB.Delete(&room)
	if result.RowsAffected == 0 {
		return errors.New("Database got an error")
	}
	return nil
}

func DeleteRoomHandler(c *fiber.Ctx) error {
	type Body struct {
		UserID uint `json:"user_id"`
	}
	roomID, roomIDerr := strconv.Atoi(c.Params("id"))
	b := new(Body)
	if err := c.BodyParser(b); err != nil || b.UserID == 0 {
		e := routes.NewErrorResponse([]string{"user_id is required"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	if roomIDerr != nil {
		e := routes.NewErrorResponse([]string{"Room not found"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	err := deleteRoomByID(roomID, b.UserID)
	if err != nil {
		e := routes.NewErrorResponse([]string{err.Error()})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	return c.SendStatus(fiber.StatusOK)
}

func joinRoom(userID uint, roomID uint, money int) (*RoomSchemaResponse, error) {
	var waitingLists []db.WaitingList
	db.DB.Where("room_id = ?", roomID).Find(&waitingLists)
	numUsers := len(waitingLists)
	if numUsers >= 9 {
		return &RoomSchemaResponse{}, errors.New("Room is full")
	}
	// user is already in waiting list
	for _, v := range waitingLists {
		if v.UserID == userID {
			return &RoomSchemaResponse{},
				errors.New("You are already in the waiting list")
		}
	}

	// check valid money (money user brings to the room must be less than or
	// equal total money of user)
	var user db.User
	db.DB.First(&user, userID)
	if money > user.Money {
		return &RoomSchemaResponse{},
			errors.New("Money is not enough")
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		waitingInfo := db.WaitingList{
			UserID:         userID,
			RoomID:         roomID,
			AvailableMoney: money,
			Ready:          false,
		}
		if err := tx.Create(&waitingInfo).Error; err != nil {
			return err
		}
		// subtract user's money in their profile
		user.Money = user.Money - money
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &RoomSchemaResponse{}, err
	}

	room, ok := getRoomById(int(roomID))
	if !ok {
		return &RoomSchemaResponse{},
			errors.New("Can not get room info")
	}
	return room, nil
}

func JoinRoomHandler(c *fiber.Ctx) error {
	type Body struct {
		UserID uint `json:"user_id"`
		Money  int  `json:"money"`
	}
	roomID, roomIDerr := strconv.Atoi(c.Params("id"))
	body := new(Body)
	if err := c.BodyParser(body); err != nil || body.UserID == 0 {
		e := routes.NewErrorResponse([]string{"user_id and money are required"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	if roomIDerr != nil {
		e := routes.NewErrorResponse([]string{"Room not found"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	roomResponse, err := joinRoom(body.UserID, uint(roomID), body.Money)
	if err != nil {
		e := routes.NewErrorResponse([]string{err.Error()})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	return c.Status(fiber.StatusOK).JSON(roomResponse)
}
