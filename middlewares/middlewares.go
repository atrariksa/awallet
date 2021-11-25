package middlewares

import (
	"fmt"
	"net/http"

	"github.com/atrariksa/awallet/configs"
	"github.com/golang-jwt/jwt"
)

func AuthMiddlewareHandler(cfg *configs.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return isAuthorized(next, cfg)
	}
}

func isAuthorized(next http.Handler, cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error Handling Token")
				}
				return []byte(cfg.JWT.Secret), nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				next.ServeHTTP(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		}
	})
}
