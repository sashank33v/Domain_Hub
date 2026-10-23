package messaging

import (
	"context"
	"domainhub/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type DomainCreatedEvent struct {
	EventType  string    `json:"event_type"`
	DomainID   int       `json:"domain_id"`
	DomainName string    `json:"domain_name"`
	Registrar  string    `json:"registrar"`
	ExpiryDate time.Time `json:"expiry_date"`
	Status     string    `json:"status"`
}

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(broker, topic string) *kafkaProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
		},
	}
}

func (p *kafkaProducer) PublishDomainCreated(ctx context.Context, domain *models.Domain) error {
	event := DomainCreatedEvent{
		EventType:  "domain.created",
		DomainID:   domain.ID,
		DomainName: domain.DomainName,
		Registrar:  domain.Registrar,
		ExpiryDate: domain.ExpiryDate,
		Status:     domain.Status,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal domain events: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   fmt.Appendf(nil, "%d", domain.ID),
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("Publish domain events: %w", err)
	}
	return nil
}

func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}
