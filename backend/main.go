package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"github.com/taiprogramer/simple-poker-game/backend/db"
	authRouter "github.com/taiprogramer/simple-poker-game/backend/routes/auth"
	roomRouter "github.com/taiprogramer/simple-poker-game/backend/routes/room"
	socketRouter "github.com/taiprogramer/simple-poker-game/backend/routes/socket"
)

func main() {
	if godotenv.Load() != nil {
		return
	}

	err := db.InitDB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	app := fiber.New()

	// non-restricted
	app.Post("/user", authRouter.SignUpHandler)
	app.Post("/auth", authRouter.SignInHandler)
	app.Get("/room", roomRouter.GetListRoomsHandler)
	app.Get("/room/:id", roomRouter.GetSpecificRoomHandler)

	// Websocket
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:user_id", websocket.New(socketRouter.SocketHandler))

	// Bearer Token is Required
	app.Use(authRouter.JWTMiddleWare())

	// restricted
	app.Get("user/:id", authRouter.GetUserHandler)
	app.Post("/room", roomRouter.CreateNewRoomHandler)
	app.Delete("/room/:id", roomRouter.DeleteRoomHandler)
	app.Post("/room/:id", roomRouter.JoinRoomHandler)
	app.Put("/room/:id", roomRouter.UpdateReadyStatusHandler)
	app.Get("/table/:id", roomRouter.GetTableHandler)

	app.Listen(":3000")
}
