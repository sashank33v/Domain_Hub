package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type kafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(broker, topic, groupID string) *kafkaConsumer {
	return &kafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *kafkaConsumer) Consume(ctx context.Context) error {
	defer c.reader.Close()

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("read kafka message: %w", err)
		}

		var event DomainCreatedEvent

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Invalid kafka event: %v", err)
			continue
		}

		log.Printf("Recieved event: type=%s domain_id=%d domain_name=%s", event.EventType, event.DomainID, event.DomainName)

	}
}
