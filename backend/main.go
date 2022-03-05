package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taiprogramer/simple-poker-game/backend/db"
	"github.com/taiprogramer/simple-poker-game/backend/utils"
)

type UserAccountSignUpBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	ErrorMessages []string `json:"error_messages"`
}

func NewErrorResponse(messages []string) ErrorResponse {
	errorResponse := ErrorResponse{}
	errorResponse.ErrorMessages = messages
	return errorResponse
}

type UserSchemaResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Money    int    `json:"money"`
}

func accountExists(user *UserAccountSignUpBody) bool {
	var u db.User
	db.DB.Where(&db.User{Username: user.Username}).First(&u)
	return u.Username != ""
}

func createAccount(body *UserAccountSignUpBody) (*UserSchemaResponse, bool) {
	hashedPassword, ok := utils.HashPassword(body.Password)
	if !ok {
		return &UserSchemaResponse{}, false
	}
	user := db.User{
		Username:       body.Username,
		HashedPassword: hashedPassword,
		Money:          0,
	}
	ok = db.DB.Create(&user).RowsAffected == 1
	return &UserSchemaResponse{
		ID:       user.ID,
		Username: user.Username,
		Money:    user.Money,
	}, ok
}

func main() {
	if !db.InitDB() {
		return
	}
	app := fiber.New()

	// create new account
	app.Post("/user", func(c *fiber.Ctx) error {
		user := new(UserAccountSignUpBody)
		if err := c.BodyParser(user); err != nil {
			e := NewErrorResponse([]string{"Unknown error"})
			return c.Status(fiber.StatusBadRequest).JSON(e)
		}

		if accountExists(user) {
			e := NewErrorResponse([]string{
				"Username already exists",
			})
			return c.Status(fiber.StatusBadRequest).JSON(e)
		}

		u, ok := createAccount(user)
		if !ok {
			e := NewErrorResponse([]string{"Unknown error"})
			return c.Status(fiber.StatusBadRequest).JSON(e)
		}

		return c.Status(fiber.StatusCreated).JSON(u)
	})

	app.Listen(":3000")
}
