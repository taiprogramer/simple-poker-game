package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taiprogramer/simple-poker-game/backend/db"
	"github.com/taiprogramer/simple-poker-game/backend/routes"
)

func main() {
	if !db.InitDB() {
		return
	}
	app := fiber.New()

	// create new account
	app.Post("/user", routes.SignUpHandler)

	app.Listen(":3000")
}
