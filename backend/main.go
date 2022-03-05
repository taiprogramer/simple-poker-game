package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID             uint `gorm:"primaryKey;autoIncrement"`
	Username       string
	HashedPassword string
	Money          int
	//WaitingLists     []*WaitingList     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BetHistories            []*BetHistory             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UsersTablesCards        []*UsersTablesCard        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UsersTablesCombinations []*UsersTablesCombination `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Action struct {
	ID           uint `gorm:"primaryKey;autoIncrement"`
	Name         string
	BetHistories []*BetHistory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Combination struct {
	ID                      uint `gorm:"primaryKey;autoIncrement"`
	Name                    string
	Score                   uint
	UsersTablesCombinations []*UsersTablesCombination `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Suit int

const (
	Diamond Suit = 0
	Heart
	Club
	Spade
)

type Card struct {
	ID     uint `gorm:"primaryKey;autoIncrement"`
	Number uint
	Suit   Suit
	Image  string
	// table has common cards
	Tables                  []*Table                  `gorm:"many2many:table_cards;"`
	UsersTablesCards        []*UsersTablesCard        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CombinationDetailsCards []*CombinationDetailsCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Room struct {
	ID       uint `gorm:"primaryKey;autoIncrement"`
	Code     string
	Playing  bool
	Private  bool
	Password string
	// setup belongs to relation
	// 1 room has 1 user as owner
	UserID       uint
	User         User
	WaitingLists []*WaitingList `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type WaitingList struct {
	// when users join room, they will appear in wating list
	UserID         uint `gorm:"primaryKey"`
	RoomID         uint
	AvailableMoney int
	Ready          bool
}

type Table struct {
	ID    uint `gorm:"primaryKey;autoIncrement"`
	Round int
	Done  bool
	Pot   int
	// table has current turn user
	UserID uint
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// table belongs to one room
	RoomID uint
	Room   Room `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// table has common cards
	Cards                   []*Card                   `gorm:"many2many:table_cards;"`
	BetHistories            []*BetHistory             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UsersTablesCards        []*UsersTablesCard        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UsersTablesCombinations []*UsersTablesCombination `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type BetHistory struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	// Bet history belongs to table, user, action (with the amount)
	TableID  uint
	Table    Table `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID   uint
	User     User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ActionID uint
	Action   Action `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount   int
}

// store cards of user on specific table
type UsersTablesCard struct {
	TableID uint
	Table   Table `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;primaryKey;"`
	UserID  uint
	User    User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;primaryKey;"`
	CardID  uint
	Card    Card `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;primaryKey;"`
}

type CombinationDetail struct {
	ID                      uint                      `gorm:"primaryKey;autoIncrement"`
	CombinationDetailsCards []*CombinationDetailsCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// store best combination of user on specific table
type UsersTablesCombination struct {
	TableID             uint
	Table               Table `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;primaryKey;"`
	UserID              uint
	User                User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;primaryKey;"`
	CombinationID       uint
	Combination         Combination `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CombinationDetailID uint
	CombinationDetail   CombinationDetail `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// store list of cards that produce a specific combination
type CombinationDetailsCard struct {
	CardID              uint
	Card                Card `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;primaryKey;"`
	CombinationDetailID uint
	CombinationDetail   CombinationDetail `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;primaryKey;"`
}

type UserAccountSignUpBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	ErrorMessages []string `json:"error_messages"`
}

type UserSchemaResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Money    int    `json:"money"`
}

func accountExists(user *UserAccountSignUpBody) bool {
	return false
}

func createAccount(user *UserAccountSignUpBody) (*UserSchemaResponse, bool) {
	u := &UserSchemaResponse{
		ID:       0,
		Username: "dummy",
		Money:    0,
	}
	return u, true
}

func main() {
	dsn := "host=localhost user=postgres password= dbname=postgres port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print("Can not connect to database")
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Action{})
	db.AutoMigrate(&Combination{})
	db.AutoMigrate(&Card{})
	db.AutoMigrate(&CombinationDetail{})
	db.AutoMigrate(&Room{})
	db.AutoMigrate(&WaitingList{})
	db.AutoMigrate(&Table{})
	db.AutoMigrate(&BetHistory{})
	db.AutoMigrate(&UsersTablesCard{})
	db.AutoMigrate(&UsersTablesCombination{})
	db.AutoMigrate(&CombinationDetailsCard{})

	app := fiber.New()

	// create new account
	app.Post("/user", func(c *fiber.Ctx) error {
		user := new(UserAccountSignUpBody)
		e := ErrorResponse{}
		if err := c.BodyParser(user); err != nil {
			e.ErrorMessages = append(e.ErrorMessages, "Unknown error")
			return c.Status(fiber.StatusBadRequest).JSON(e)
		}

		if accountExists(user) {
			e.ErrorMessages = append(e.ErrorMessages, "Username already exists")
			return c.Status(fiber.StatusBadRequest).JSON(e)
		}

		u, ok := createAccount(user)
		if !ok {
			e.ErrorMessages = append(e.ErrorMessages, "Unknown error")
			return c.Status(fiber.StatusBadRequest).JSON(e)
		}

		return c.Status(fiber.StatusCreated).JSON(u)
	})

	app.Listen(":3000")
}
