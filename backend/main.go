package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taiprogramer/simple-poker-game/backend/db"
)

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
	if !db.InitDB() {
		return
	}
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
