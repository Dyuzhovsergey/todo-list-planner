// Package auth for authentication
package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"github.com/Dyuzhovsergey/todo-list-planner/pkg/config"
)

// AuthMiddleware оборачивает хэндлеры, требующие аутентификации
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.Cfg.Password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return config.Cfg.JWTKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}

}
