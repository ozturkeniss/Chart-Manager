package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

type Publisher struct {
	producer sarama.SyncProducer
}

type StockUpdateEvent struct {
	ItemID    string `json:"item_id"`
	NewStock  int    `json:"new_stock"`
	UserID    uint32 `json:"user_id"`
	EventType string `json:"event_type"`
}

type ItemDeleteEvent struct {
	ItemID    string `json:"item_id"`
	UserID    uint32 `json:"user_id"`
	EventType string `json:"event_type"`
}

func NewPublisher(brokers []string) (*Publisher, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	return &Publisher{producer: producer}, nil
}

func (p *Publisher) PublishStockUpdate(event *StockUpdateEvent) error {
	event.EventType = "stock_update"

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "stock_updates",
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Stock update event published to partition %d at offset %d", partition, offset)
	return nil
}

func (p *Publisher) PublishItemDelete(event *ItemDeleteEvent) error {
	event.EventType = "item_delete"

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "item_deletes",
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Item delete event published to partition %d at offset %d", partition, offset)
	return nil
}

func (p *Publisher) Close() error {
	return p.producer.Close()
}
