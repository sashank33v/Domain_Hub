package handlers

import (
	"domainhub/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatHandlerBadRequests(t *testing.T) {
	h := NewChatHandler(service.NewAIService(nil, "keu", ""))

	tests := []struct {
		name string
		body string
	}{
		{"invalid json", `{bad`},
		{"empty message", `{"message":"  "}`},
		{"message too long", `{"message":"}` + strings.Repeat("a", 501) + `"}`},
		{"body too large", `{"message":"}` + strings.Repeat("a", 5000) + `"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			h.Chat(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected %d, got %d", http.StatusBadRequest, rec.Code)
			}
		})
	}
}
