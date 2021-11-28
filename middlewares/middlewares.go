package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/atrariksa/awallet/configs"
	"github.com/atrariksa/awallet/models"
	"github.com/golang-jwt/jwt"
)

func AuthMiddlewareHandler(cfg *configs.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return isAuthorized(next, cfg)
	}
}

func isAuthorized(next http.Handler, cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {

			token, err := jwt.ParseWithClaims(r.Header["Authorization"][0], &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error Handling Token")
				}
				return []byte(cfg.JWT.Secret), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
				return
			}

			if token == nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
				return
			}

			if claims, ok := token.Claims.(*models.JwtClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "token", claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
				return
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}
	})
}
