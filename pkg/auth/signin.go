package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dyuzhovsergey/todo-list-planner/pkg/config"
	"github.com/golang-jwt/jwt/v4"
)

type SigninRequest struct {
	Password string `json:"password"`
}

type SigninResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var req SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if config.Cfg.Password == "" || req.Password != config.Cfg.Password {
		json.NewEncoder(w).Encode(SigninResponse{Error: "invalid passord"})
		return
	}

	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(config.Cfg.JWTKey)
	if err != nil {
		http.Error(w, "cannot generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
	})
	json.NewEncoder(w).Encode(SigninResponse{Token: tokenStr})

}
