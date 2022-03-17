package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/taiprogramer/simple-poker-game/backend/db"
	authRouter "github.com/taiprogramer/simple-poker-game/backend/routes/auth"
	roomRouter "github.com/taiprogramer/simple-poker-game/backend/routes/room"
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
	app.Post("/user", authRouter.SignUpHandler)
	app.Post("/auth", authRouter.SignInHandler)
	app.Get("/room", roomRouter.GetListRoomsHandler)
	app.Get("/room/:id", roomRouter.GetSpecificRoomHandler)

	// Bearer Token is Required
	app.Use(authRouter.JWTMiddleWare())

	// restricted
	app.Get("user/:id", authRouter.GetUserHandler)
	app.Post("/room", roomRouter.CreateNewRoomHandler)
	app.Delete("/room/:id", roomRouter.DeleteRoomHandler)

	app.Listen(":3000")
}
