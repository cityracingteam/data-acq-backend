package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	ExtraClaims
}

type ExtraClaims struct {
	KeyID string `json:"kid,omitempty"`
}

type ExpiryOffset struct {
	years  int
	months int
	days   int
}
