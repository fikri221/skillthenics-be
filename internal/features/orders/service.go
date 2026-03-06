package orders

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/ksuid"
)

type Service interface {
	GetOrder(ctx context.Context, id string) (Order, error)
	CreateOrder(ctx context.Context, customerName string, items []OrderItem) (string, error)
}

type svc struct {
	repo        Repository
	kafkaWriter *kafka.Writer
}

func NewService(repo Repository, kafkaWriter *kafka.Writer) Service {
	return &svc{repo: repo, kafkaWriter: kafkaWriter}
}

func (s *svc) GetOrder(ctx context.Context, id string) (Order, error) {
	return s.repo.GetOrderWithItems(ctx, id)
}

func (s *svc) CreateOrder(ctx context.Context, customerName string, items []OrderItem) (string, error) {
	orderID := ksuid.New().String()

	var total float64
	for i := range items {
		total += float64(items[i].Quantity) * items[i].Price
	}

	order := Order{
		ID:           orderID,
		CustomerName: customerName,
		TotalAmount:  total,
		Items:        items,
	}

	err := s.repo.WithTx(ctx, func(txRepo Repository) error {
		return txRepo.CreateOrder(ctx, order)
	})

	if err != nil {
		return "", err
	}

	// Emit Kafka Event
	if s.kafkaWriter != nil {
		eventData, _ := json.Marshal(order)
		err := s.kafkaWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(orderID),
			Value: eventData,
		})
		if err != nil {
			fmt.Printf("failed to emit kafka event: %v\n", err)
		}
	}

	return orderID, nil
}
