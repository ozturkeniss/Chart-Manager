package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	handler  EventHandler
}

type EventHandler interface {
	HandleStockUpdate(event *StockUpdateEvent) error
	HandleItemDelete(event *ItemDeleteEvent) error
}

func NewConsumer(brokers []string, groupID string, handler EventHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	return &Consumer{
		consumer: consumer,
		handler:  handler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context, topics []string) error {
	for {
		err := c.consumer.Consume(ctx, topics, c)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message topic:%q partition:%d offset:%d\n", message.Topic, message.Partition, message.Offset)

		switch message.Topic {
		case "stock_updates":
			var event StockUpdateEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("Failed to unmarshal stock update event: %v", err)
				continue
			}
			if err := c.handler.HandleStockUpdate(&event); err != nil {
				log.Printf("Failed to handle stock update event: %v", err)
			}

		case "item_deletes":
			var event ItemDeleteEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("Failed to unmarshal item delete event: %v", err)
				continue
			}
			if err := c.handler.HandleItemDelete(&event); err != nil {
				log.Printf("Failed to handle item delete event: %v", err)
			}
		}

		session.MarkMessage(message, "")
	}

	return nil
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
