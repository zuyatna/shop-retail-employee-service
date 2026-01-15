package adapterhttp

import (
	"encoding/json"
	"net/http"

	"github.com/zuyatna/shop-retail-employee-service/internal/dto/auth"
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
)

type AuthHandler struct {
	usecase *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: uc,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	token, err := h.usecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		WriteErrorJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, auth.LoginResponse{Token: token}, "login successful")
}
