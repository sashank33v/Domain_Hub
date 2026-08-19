package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey contextKey = "user"

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(
					w,
					"Missing authorization header",
					http.StatusUnauthorized,
				)
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(
					w,
					"Invalid authorization header",
					http.StatusUnauthorized,
				)
				return
			}
			tokenString := parts[1]

			token, err := jwt.Parse(
				tokenString,
				func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, jwt.ErrTokenSignatureInvalid
					}
					return []byte(jwtSecret), nil
				},
			)
			if err != nil || !token.Valid {
				http.Error(
					w,
					"Invalid or expired token",
					http.StatusUnauthorized,
				)
				return
			}
			next.ServeHTTP(w, r)

		})
	}
}
