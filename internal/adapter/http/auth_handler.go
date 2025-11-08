package http

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/zuyatna/shop-retail-employee-service/internal/model"
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
)

type AuthHandler struct {
	auth *usecase.AuthUsecase
}

func NewAuthHandler(auth *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	token, _, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadRequest):
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password are required"})
		case errors.Is(err, domain.ErrNotFound), errors.Is(err, domain.ErrDeleted), errors.Is(err, usecase.ErrInvalidCredentials):
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	writeJSON(w, http.StatusOK, loginResponse{Token: token})
}
