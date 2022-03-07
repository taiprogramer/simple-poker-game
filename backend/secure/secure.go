package secure

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, bool) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(b), err == nil
}

func ComparePassword(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}

func GenerateToken(s string) string {
	return ""
}
