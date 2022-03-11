package secure

import (
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

const expectedButActual = "Expected: %v but Actual %v"
const incorrectPass = "incorrect_pass"

// bcrypt hash of "pass"
const h = "$2y$12$Tf9QJXpknE8JiC70dQBYM.83aHMG/QCJhfbKxwcc23/Arv9TSe/om"
const plain = "pass"

// jwt: {usr: "trump"}
const tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c3IiOiJ0cnVtcCJ9.20PTxtLpqDSPBIUsYpegw6Zfon-QdeYDcg3_COMdoZE"
const claimUser = "trump"

func TestHashPassword(t *testing.T) {
	h, _ := HashPassword(plain)
	l := len(h)
	if l != 60 {
		t.Fatalf(expectedButActual, 60, l)
	}
}

func TestComparePassword(t *testing.T) {
	if !ComparePassword(plain, h) {
		t.Fatalf(expectedButActual, true, false)
	}
	if ComparePassword(incorrectPass, h) {
		t.Fatalf(expectedButActual, false, true)
	}
}

func TestGenerateToken(t *testing.T) {
	tokenString, ok := GenerateToken("simple")
	if !ok {
		t.Fatalf(expectedButActual, true, ok)
	}

	// jwt format must have 2 dots: xxxxxxxx.xxxxxxxx.xxxxxxxx
	// Header, Payload, Signature
	dots := strings.Count(tokenString, ".")
	if dots != 2 {
		t.Fatalf(expectedButActual, 2, dots)
	}
}

func TestGetJWTClaim(t *testing.T) {
	token, _ := jwt.Parse(tokenString, nil)
	usr := GetJWTClaim("usr", token)
	if usr != claimUser {
		t.Fatalf(expectedButActual, claimUser, usr)
	}
}
