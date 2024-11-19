package secure

import (
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

const expectedButActual = "Expected: %v but Actual %v"
const incorrectPass = "incorrect_pass"

func CheckExpect(expect, actual interface{}, t *testing.T) {
	if expect != actual {
		t.Fatalf(expectedButActual, expect, actual)
	}
}

// bcrypt hash of "pass"
const h = "$2y$12$Tf9QJXpknE8JiC70dQBYM.83aHMG/QCJhfbKxwcc23/Arv9TSe/om"
const plain = "pass"

// jwt: {usr: "trump"}
const tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c3IiOiJ0cnVtcCJ9.20PTxtLpqDSPBIUsYpegw6Zfon-QdeYDcg3_COMdoZE"
const claimUser = "trump"

func TestHashPassword(t *testing.T) {
	h, _ := HashPassword(plain)
	l := len(h)
	CheckExpect(60, l, t)
}

func TestComparePassword(t *testing.T) {
	CheckExpect(true, ComparePassword(plain, h), t)
	CheckExpect(false, ComparePassword(incorrectPass, h), t)
}

func TestGenerateToken(t *testing.T) {
	tokenString, ok := GenerateToken("simple")
	CheckExpect(true, ok, t)

	// jwt format must have 2 dots: xxxxxxxx.xxxxxxxx.xxxxxxxx
	// Header, Payload, Signature
	dots := strings.Count(tokenString, ".")
	CheckExpect(2, dots, t)
}

func TestGetJWTClaim(t *testing.T) {
	token, _ := jwt.Parse(tokenString, nil)
	usr := GetJWTClaim("usr", token)
	CheckExpect(claimUser, usr, t)
}
