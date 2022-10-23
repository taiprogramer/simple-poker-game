package room

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/taiprogramer/simple-poker-game/backend/db"
	"github.com/taiprogramer/simple-poker-game/backend/repo/bet_histories"
	tableRepo "github.com/taiprogramer/simple-poker-game/backend/repo/table"
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

type CardSchemaResponse struct {
	ID     uint   `json:"id"`
	Number int    `json:"number"`
	Suit   int    `json:"suit"`
	Image  string `json:"image"`
}

type CombinationSchemaResponse struct {
	ID                  uint                 `json:"id"`
	Name                string               `json:"name"`
	SelectedCommonCards []CardSchemaResponse `json:"selected_common_cards"`
}

type TurnSchemaResponse struct {
	UserID uint `json:"user_id"`
}

type ActionSchemaResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type BetSchemaResponse struct {
	UserID uint                 `json:"user_id"`
	Action ActionSchemaResponse `json:"action"`
}

type AvailableActionSchemaResponse struct {
	CanFold  bool `json:"can_fold"`
	CanCheck bool `json:"can_check"`
	CanRaise bool `json:"can_raise"`
	CanCall  bool `json:"can_call"`
}

type PlayerSchemaResponse struct {
	ID             uint `json:"id"`
	AvailableMoney int  `json:"available_money"`
}

type ResultSchemaResponse struct {
	UserID uint                 `json:"user_id"`
	Cards  []CardSchemaResponse `json:"cards"`
}

type WinningSchemaResponse struct {
	UserID          uint                      `json:"user_id"`
	WinningAmount   int                       `json:"winning_amount"`
	BestCombination CombinationSchemaResponse `json:"best_combination"`
}

type TableSchemaResponse struct {
	ID              uint                          `json:"id"`
	Round           int                           `json:"round"`
	Done            bool                          `json:"done"`
	Pot             int                           `json:"pot"`
	CommonCards     []CardSchemaResponse          `json:"common_cards"`
	OwnCards        []CardSchemaResponse          `json:"own_cards"`
	BestCombination CombinationSchemaResponse     `json:"best_combination"`
	CurrentTurn     TurnSchemaResponse            `json:"current_turn"`
	LatestBet       BetSchemaResponse             `json:"latest_bet"`
	AvailableAction AvailableActionSchemaResponse `json:"available_action"`
	Players         []PlayerSchemaResponse        `json:"players"`
	Results         []ResultSchemaResponse        `json:"results"`
	Winers          []WinningSchemaResponse       `json:"winers"`
}

type RoomBody struct {
	UserID   uint   `json:"user_id"`
	Money    int    `json:"money"`
	Password string `json:"password"`
}

func isAuthenticatedUser(userID uint, username interface{}) bool {
	var user db.User
	db.DB.First(&user, userID)
	return user.Username == username
}

// Inspired by YuMy
// Note: I know this implementation has a draw back that means collision can be
// occurred. Maybe this implementation will be replaced in the future.
func getNextUniqueCode() string {
	rand.Seed(time.Now().UnixNano())
	var code string
	charList := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := 0; i < 4; i++ {
		index := rand.Int() % len(charList)
		code += charList[index]
	}

	return code
}

func createNewRoom(userID uint, private bool, password string, money int) (*RoomSchemaResponse, bool) {
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
		AvailableMoney: money,
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
	response, ok := createNewRoom(body.UserID, private, body.Password, body.Money)

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

func updateReadyStatus(userID uint, roomID int, ready bool) (*RoomSchemaResponse, error) {
	var waitingList db.WaitingList
	result := db.DB.Where(
		"room_id = ? AND user_id = ?",
		roomID,
		userID).First(&waitingList)

	if result.RowsAffected == 0 {
		return nil, errors.New("Please join room first")
	}

	waitingList.Ready = ready
	result = db.DB.Save(&waitingList)
	if result.RowsAffected == 0 {
		return nil, errors.New("Can not update ready status")
	}

	room, ok := getRoomById(int(roomID))
	if !ok {
		return nil, errors.New("Can not get room info")
	}
	return room, nil
}

func UpdateReadyStatusHandler(c *fiber.Ctx) error {
	type Body struct {
		UserID uint `json:"user_id"`
		Ready  bool `json:"ready"`
	}

	roomID, roomIDerr := strconv.Atoi(c.Params("id"))
	body := new(Body)
	if err := c.BodyParser(body); err != nil || body.UserID == 0 {
		e := routes.NewErrorResponse(
			[]string{"user_id and ready status are required"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	if roomIDerr != nil {
		e := routes.NewErrorResponse([]string{"Room not found"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	roomResponse, err := updateReadyStatus(body.UserID, roomID, body.Ready)
	if err != nil {
		e := routes.NewErrorResponse([]string{err.Error()})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	return c.Status(fiber.StatusOK).JSON(roomResponse)
}

func cardsToCardsResponse(cards []*db.Card) []CardSchemaResponse {
	cardsSchema := make([]CardSchemaResponse, 0)
	for _, v := range cards {
		c := CardSchemaResponse{
			ID:     v.ID,
			Number: int(v.Number),
			Suit:   int(v.Suit),
			Image:  v.Image,
		}
		cardsSchema = append(cardsSchema, c)
	}
	return cardsSchema
}

func userTableCardsToCards(value []db.UsersTablesCard) []*db.Card {
	var cards []*db.Card
	for i := 0; i < len(value); i++ {
		cards = append(cards, &value[i].Card)
	}
	return cards
}

func combinationDetailCardsToCards(value []db.CombinationDetailsCard) []*db.Card {
	cards := make([]*db.Card, 0)
	for _, v := range value {
		cards = append(cards, &v.Card)
	}
	return cards
}

func waitingListToPlayers(value []db.WaitingList) []PlayerSchemaResponse {
	players := make([]PlayerSchemaResponse, 0)
	for _, v := range value {
		players = append(players, PlayerSchemaResponse{
			ID:             v.UserID,
			AvailableMoney: v.AvailableMoney,
		})
	}
	return players
}

func waitingListToResults(value []db.WaitingList, tableID uint) []ResultSchemaResponse {
	results := make([]ResultSchemaResponse, 0)
	for _, v := range value {
		var cardsOfUser []db.UsersTablesCard
		db.DB.Preload("Card").Where(&db.UsersTablesCard{TableID: tableID, UserID: v.UserID}).Find(&cardsOfUser)
		results = append(results, ResultSchemaResponse{
			UserID: v.UserID,
			Cards:  cardsToCardsResponse(userTableCardsToCards(cardsOfUser)),
		})
	}
	return results
}

func getTableInfo(tableID int, userID uint) (*TableSchemaResponse, error) {
	var table db.Table
	result := db.DB.Preload("Cards").First(&table, tableID)
	if result.RowsAffected == 0 {
		return nil, errors.New("Table not found")
	}

	var ownCards []db.UsersTablesCard
	db.DB.Preload("Card").Where(&db.UsersTablesCard{TableID: table.ID, UserID: userID}).Find(&ownCards)

	var bestCombination db.UsersTablesCombination
	db.DB.Preload("Combination").Where(&db.UsersTablesCombination{TableID: table.ID, UserID: userID}).First(&bestCombination)
	var combinationDetailCards []db.CombinationDetailsCard
	db.DB.Preload("Card").Where(&db.CombinationDetailsCard{CombinationDetailID: bestCombination.CombinationDetailID}).Find(&combinationDetailCards)
	var betHistory db.BetHistory
	db.DB.Preload("Action").Where(&db.BetHistory{TableID: table.ID}).Last(&betHistory)
	var waitingList []db.WaitingList
	db.DB.Where(&db.WaitingList{Ready: true, RoomID: table.RoomID}).Find(&waitingList)

	response := &TableSchemaResponse{
		ID:    table.ID,
		Round: table.Round,
		Done:  table.Done,
		Pot:   table.Pot,
		CurrentTurn: TurnSchemaResponse{
			UserID: table.UserID,
		},
		CommonCards: cardsToCardsResponse(table.Cards),
		OwnCards:    cardsToCardsResponse(userTableCardsToCards(ownCards)),
		BestCombination: CombinationSchemaResponse{
			ID:                  bestCombination.CombinationID,
			Name:                bestCombination.Combination.Name,
			SelectedCommonCards: cardsToCardsResponse(combinationDetailCardsToCards(combinationDetailCards)),
		},
		LatestBet: BetSchemaResponse{
			UserID: betHistory.UserID,
			Action: ActionSchemaResponse{
				ID:     betHistory.ActionID,
				Name:   betHistory.Action.Name,
				Amount: betHistory.Amount,
			},
		},
		Players: waitingListToPlayers(waitingList),
		// only available when the game is done.
		Results: waitingListToResults(waitingList, table.ID),
	}

	return response, nil
}

func GetTableHandler(c *fiber.Ctx) error {
	type Query struct {
		UserID uint
	}
	var query Query
	tableID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		e := routes.NewErrorResponse(
			[]string{"Table not found"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	c.QueryParser(&query)
	if query.UserID == 0 {
		e := routes.NewErrorResponse(
			[]string{"Please provide your user id"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	response, err := getTableInfo(tableID, query.UserID)
	if err != nil {
		e := routes.NewErrorResponse(
			[]string{err.Error()})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func writeBetHistory(tableID, userID int, round int, actionName string, amount int) {
	actionID := bet_histories.FindActionIDByName(actionName)
	bet_histories.WriteBetHistory(tableID, userID, actionID, amount, round)
}

type ActionData struct {
	UserID int    `json:"user_id"`
	Action string `json:"action"`
	Amount int    `json:"amount"`
}

func performAction(actionData ActionData, tableID int) {
	table := tableRepo.GetTableByID(tableID)
	betHistory := bet_histories.GetLatestByTableID(tableID)
	totalAmountPreviousBet := bet_histories.GetTotalAmountByRoundAndUserID(int(table.ID), table.Round, int(betHistory.UserID))
	totalAmount := bet_histories.GetTotalAmountByRoundAndUserID(int(table.ID), table.Round, actionData.UserID)
	if strings.Compare(actionData.Action, "call") == 0 {
		writeBetHistory(tableID, actionData.UserID, table.Round,
			actionData.Action, totalAmountPreviousBet-totalAmount)
	}
	if strings.Compare(actionData.Action, "raise") == 0 {
		writeBetHistory(tableID, actionData.UserID, table.Round,
			actionData.Action, (totalAmountPreviousBet-totalAmount)+
				actionData.Amount)
	}
}

func PerformActionHandler(c *fiber.Ctx) error {
	tableID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		e := routes.NewErrorResponse([]string{"Please supply a table id!"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	var actionData ActionData
	c.BodyParser(&actionData)

	performAction(actionData, tableID)

	// just return the old info of the table, delegate the job update table status
	// to the socket logic (simplify the overall cost for now)
	// it's better to have the logic update table status in here for data consistency,
	// to do this, we need to re-design the database structure,
	// it's risky to do at this stage, so simplify for now.
	response, err := getTableInfo(tableID, uint(actionData.UserID))
	if err != nil {
		e := routes.NewErrorResponse(
			[]string{err.Error()})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
