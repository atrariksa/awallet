package models

import "github.com/golang-jwt/jwt"

type JwtClaims struct {
	ID       string `json:"id"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
