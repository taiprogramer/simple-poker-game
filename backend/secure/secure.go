package secure

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, bool) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(b), err == nil
}
