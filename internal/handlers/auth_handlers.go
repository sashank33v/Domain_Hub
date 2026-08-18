package handlers

import (
	"domainhub/internal/dto"
	"domainhub/internal/service"
	"domainhub/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	err = h.service.Register(&req)
	if err != nil {
		utils.SendErrorResponse(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	utils.SendResponse(
		w,
		http.StatusCreated,
		"User registered successfully",
		nil,
		nil,
	)
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	token, err := h.service.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.SendErrorResponse(
				w,
				http.StatusUnauthorized,
				"Invalid email or password",
			)
			return
		}
		utils.SendErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to login",
		)
		return
	}

	utils.SendResponse(
		w,
		http.StatusOK,
		"Login Successful",
		map[string]string{
			"token": token,
		},
		nil,
	)
}
