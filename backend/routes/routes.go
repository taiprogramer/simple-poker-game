package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/taiprogramer/simple-poker-game/backend/db"
	"github.com/taiprogramer/simple-poker-game/backend/secure"
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

type SignInSuccessResponse struct {
	AccessToken string `json:"access_token"`
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
		Money:          0,
	}
	ok = db.DB.Create(&user).RowsAffected == 1
	return &UserSchemaResponse{
		ID:       user.ID,
		Username: user.Username,
		Money:    user.Money,
	}, ok
}

func userAndPasswordCorrect(body *UserAccountSignUpBody) bool {
	var user db.User
	result := db.DB.Where("username = ?", body.Username).First(&user)
	if result.RowsAffected == 0 {
		return false
	}
	if !secure.ComparePassword(body.Password, user.HashedPassword) {
		return false
	}
	return true
}

func SignUpHandler(c *fiber.Ctx) error {
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
}

func SignInHandler(c *fiber.Ctx) error {
	user := new(UserAccountSignUpBody)
	if err := c.BodyParser(user); err != nil {
		e := NewErrorResponse([]string{"Unknown error"})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	if !userAndPasswordCorrect(user) {
		e := NewErrorResponse([]string{
			"Incorrect username or password",
		})
		return c.Status(fiber.StatusBadRequest).JSON(e)
	}

	tokenString, ok := secure.GenerateToken(user.Username)
	if !ok {
		e := NewErrorResponse([]string{
			"Can not generate jwt token",
		})
		return c.Status(fiber.StatusInternalServerError).JSON(e)
	}
	response := SignInSuccessResponse{
		AccessToken: tokenString,
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}
