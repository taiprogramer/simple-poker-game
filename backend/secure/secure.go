package secure

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, bool) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(b), err == nil
}

func ComparePassword(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}

func GenerateToken(s string) (string, bool) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": s,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	})
	hmacSecret := os.Getenv("HMAC_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(hmacSecret))
	return tokenString, err == nil
}

func GetJWTClaim(claim string, token interface{}) interface{} {
	_token := token.(*jwt.Token)
	claims := _token.Claims.(jwt.MapClaims)
	return claims[claim]
}
