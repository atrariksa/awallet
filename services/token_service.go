package services

import (
	"time"

	"github.com/atrariksa/awallet/configs"
	"github.com/atrariksa/awallet/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenService struct {
	cfg *configs.Config
}

func NewTokenService(cfg *configs.Config) TokenService {
	return TokenService{
		cfg: cfg,
	}
}

type ITokenService interface {
	CreateToken(userID uint) (string, error)
}

func (ts *TokenService) CreateToken(userID uint) (signedString string, err error) {
	uuidStr := uuid.New().String()
	claims := &models.JwtClaims{
		ID:     uuidStr,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err = token.SignedString([]byte(ts.cfg.JWT.Secret))
	return
}
