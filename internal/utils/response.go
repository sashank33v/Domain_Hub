package utils

import (
	"domainhub/internal/models"
	"encoding/json"
	"net/http"
)

func SendResponse(
	w http.ResponseWriter,
	statusCode int,
	message string,
	data interface{},
	pagination *models.Pagination,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.Response{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}

	json.NewEncoder(w).Encode(response)
}
func SendErrorResponse(
	w http.ResponseWriter,
	statusCode int,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.Response{
		Success: false,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}
