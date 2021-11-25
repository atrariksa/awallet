package models

import "github.com/golang-jwt/jwt"

type JwtClaims struct {
	ID     string `json:"id"`
	UserID uint   `json:"username"`
	jwt.StandardClaims
}
