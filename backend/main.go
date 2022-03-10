package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/taiprogramer/simple-poker-game/backend/db"
	"github.com/taiprogramer/simple-poker-game/backend/routes"
)

func main() {
	if godotenv.Load() != nil {
		return
	}

	if !db.InitDB() {
		return
	}
	app := fiber.New()

	// non-restricted
	app.Post("/user", routes.SignUpHandler)
	app.Post("/auth", routes.SignInHandler)

	// Bearer Token is Required
	app.Use(routes.JWTMiddleWare())

	// restricted

	app.Listen(":3000")
}
