package util

import (
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
