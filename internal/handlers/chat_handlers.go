package handlers

import (
	"domainhub/internal/service"
	"domainhub/internal/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type ChatRequest struct {
	Message string             `json:"message"`
	History []service.ChatTurn `json:"history"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

type ChatHandler struct {
	service *service.AIService
}

func NewChatHandler(service *service.AIService) *ChatHandler {
	return &ChatHandler{
		service: service,
	}
}

func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	r.Body = http.MaxBytesReader(w, r.Body, 32*1024)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Message) == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Message is required")
		return
	}

	reply, err := h.service.Chat(r.Context(), req.Message, req.History)
	if err != nil {
		if errors.Is(err, service.ErrMessageTooLong) {
			utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Printf("Chat Failed: %v\n", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to process chat message")
		return
	}
	utils.SendResponse(w, http.StatusOK, "Chat response generated", ChatResponse{Reply: reply}, nil)
}
