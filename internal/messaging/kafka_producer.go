package messaging

import (
	"context"
	"domainhub/internal/models"
	"encoding/json"
	"fmt"
	"log"
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
type DomainUpdatedEvent struct {
	EventType  string    `json:"event_type"`
	DomainID   int       `json:"domain_id"`
	DomainName string    `json:"domain_name"`
	Registrar  string    `json:"registrar"`
	ExpiryDate time.Time `json:"expiry_date"`
	Status     string    `json:"status"`
}
type DomainDeletedEvent struct {
	EventType  string `json:"event_type"`
	DomainID   int    `json:"domain_id"`
	DomainName string `json:"domain_name"`
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
func (p *kafkaProducer) writewithRetry(ctx context.Context, msg kafka.Message) error {
	var err error

	for attempt := 1; attempt <= 3; attempt++ {
		log.Printf("Kafka publish attempt [%d/3]", attempt)

		err = p.writer.WriteMessages(ctx, msg)
		if err == nil {
			log.Printf("Kafka message published successfully")
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if attempt < 3 {
			delay := time.Duration(attempt) * 500 * time.Millisecond

			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}
	return fmt.Errorf("kafka publish failed after 3 attempts: %w", err)
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

	msg := kafka.Message{
		Key:   fmt.Appendf(nil, "%d", domain.ID),
		Value: payload,
	}
	if err := p.writewithRetry(ctx, msg); err != nil {
		return fmt.Errorf("publish domain events: %w", err)
	}
	return nil
}

func (p *kafkaProducer) PublishDomainUpdated(ctx context.Context, domain *models.Domain) error {
	event := DomainUpdatedEvent{
		EventType:  "domain.updated",
		DomainID:   domain.ID,
		DomainName: domain.DomainName,
		Registrar:  domain.Registrar,
		ExpiryDate: domain.ExpiryDate,
		Status:     domain.Status,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal domain updated event: %w", err)
	}

	msg := kafka.Message{
		Key:   fmt.Appendf(nil, "%d", domain.ID),
		Value: payload,
	}
	if err := p.writewithRetry(ctx, msg); err != nil {
		return fmt.Errorf("Publish domain updated event: %w", err)
	}
	return nil
}
func (p *kafkaProducer) PublishDomainDeleted(ctx context.Context, domainID int, domainName string) error {
	event := DomainDeletedEvent{
		EventType:  "domain.deleted",
		DomainID:   domainID,
		DomainName: domainName,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal domain deleted event: %w", err)
	}

	msg := kafka.Message{
		Key:   fmt.Appendf(nil, "%d", domainID),
		Value: payload,
	}
	if err := p.writewithRetry(ctx, msg); err != nil {
		return fmt.Errorf("Publish domain updated event: %w", err)
	}
	return nil
}
func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}
