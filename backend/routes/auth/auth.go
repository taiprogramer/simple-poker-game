package auth

import (
	"os"
	"strconv"

	"backend/constants"
	"backend/db"
	"backend/routes"
	"backend/secure"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserAccountSignUpBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSchemaResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Money    int    `json:"money"`
}

type SignInSuccessResponse struct {
	AccessToken string `json:"access_token"`
	UserID      uint   `json:"user_id"`
}

// Bearer Token Authorization Middleware (JWT)
func JWTMiddleWare() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("HMAC_SECRET_KEY"))},
		ContextKey: "token",
	})
}

func accountExists(user *UserAccountSignUpBody) bool {
	var u db.User
	db.DB.Where(&db.User{Username: user.Username}).First(&u)
	return u.Username != ""
}

func createAccount(body *UserAccountSignUpBody) (*UserSchemaResponse, bool) {
	hashedPassword, ok := secure.HashPassword(body.Password)
	if !ok {
		return &UserSchemaResponse{}, false
	}
	user := db.User{
		Username:       body.Username,
		HashedPassword: hashedPassword,
		Money:          constants.DEFAULT_MONEY_FOR_NEW_USER,
	}
	ok = db.DB.Create(&user).RowsAffected == 1
	return &UserSchemaResponse{
		ID:       user.ID,
		Username: user.Username,
		Money:    user.Money,
	}, ok
}

// return user_id if username & password are correct
func userAndPasswordCorrect(body *UserAccountSignUpBody) uint {
	var user db.User
	result := db.DB.Where("username = ?", body.Username).First(&user)
	if result.RowsAffected == 0 {
		return 0
	}
	if !secure.ComparePassword(body.Password, user.HashedPassword) {
		return 0
	}
	return user.ID
}

func getUserById(id int) (*db.User, bool) {
	var user db.User
	result := db.DB.First(&user, id)
	if result.RowsAffected == 0 {
		return nil, false
	}
	return &user, true
}

func SignUpHandler(c *fiber.Ctx) error {
	user := new(UserAccountSignUpBody)
	if err := c.BodyParser(user); err != nil {
		e := routes.NewErrorResponse([]string{"Can not parse body, try to include Content-Type: application/json in your request's header"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	if accountExists(user) {
		e := routes.NewErrorResponse([]string{
			"Username already exists",
		})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	u, ok := createAccount(user)
	if !ok {
		e := routes.NewErrorResponse([]string{"Can not access database"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	return c.Status(fiber.StatusCreated).JSON(u)
}

func SignInHandler(c *fiber.Ctx) error {
	user := new(UserAccountSignUpBody)
	if err := c.BodyParser(user); err != nil {
		e := routes.NewErrorResponse([]string{"Unknown error"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	userID := userAndPasswordCorrect(user)

	if userID == 0 {
		e := routes.NewErrorResponse([]string{
			"Incorrect username or password",
		})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	tokenString, ok := secure.GenerateToken(user.Username)
	if !ok {
		e := routes.NewErrorResponse([]string{
			"Can not generate jwt token",
		})
		return c.Status(fiber.StatusInternalServerError).JSON(e)
	}
	response := SignInSuccessResponse{
		AccessToken: tokenString,
		UserID:      userID,
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

func GetUserHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		e := routes.NewErrorResponse([]string{"Please supply your user id!"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	// get user in db
	user, ok := getUserById(id)
	if !ok {
		e := routes.NewErrorResponse([]string{"User does not exist."})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	// check correct owner
	token := c.Locals("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	username := claims["usr"]

	if username != user.Username {
		e := routes.NewErrorResponse([]string{
			"Please supply your user id! Not id of other. Are you hacker?",
		})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	response := UserSchemaResponse{
		ID:       user.ID,
		Username: user.Username,
		Money:    user.Money,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
