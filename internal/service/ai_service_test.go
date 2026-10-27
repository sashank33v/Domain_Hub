package service

import (
	"context"
	"domainhub/internal/models"
	"errors"
	"strings"
	"testing"
	"time"
)

type fakeLister struct {
	domains []models.Domain
}

func (f *fakeLister) GetAll(limit, offset int, status, registrar string) ([]models.Domain, int, error) {
	return f.domains, len(f.domains), nil
}

func TestBuildDomainContext(t *testing.T) {
	today := dateOnly(time.Now())
	s := NewAIService(&fakeLister{
		domains: []models.Domain{
			{ID: 1, DomainName: "soon.com", Registrar: "GoDaddy", ExpiryDate: today.AddDate(0, 0, 10), Status: "ACTIVE"},
			{ID: 2, DomainName: "old.org", Registrar: "Namecheap", ExpiryDate: today.AddDate(0, 0, -5), Status: "EXPIRED"},
		}}, "key", "")

	got, err := s.buildDomainContext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []string{
		"Today's date: " + today.Format("2006-01-02"),
		"soon.com | GoDaddy",
		"| ACTIVE | 10",
		"old.org | Namecheap",
		"| EXPIRED | -5",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("context missing %q\n%s", want, got)
		}
	}
}
func TestChatValidation(t *testing.T) {
	ctx := context.Background()

	s := NewAIService(&fakeLister{}, "key", "")

	if _, err := s.Chat(ctx, "   ", nil); err == nil {
		t.Error("expected error for empty message")
	}

	long := strings.Repeat("a", maxMessageRunes+1)
	if _, err := s.Chat(ctx, long, nil); !errors.Is(err, ErrMessageTooLong) {
		t.Errorf("expected ErrMessageTooLong, got %v", err)
	}

	noKey := NewAIService(&fakeLister{}, "", "")
	if _, err := noKey.Chat(ctx, "hello", nil); err == nil {
		t.Error("expected error when API key is missing")
	}
}
