package products

import (
	"context"
	"fmt"
	"nds-go-starter/internal/json"
	"net/http"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/ksuid"
)

var (
	ErrProductNameExists = json.NewError(http.StatusConflict, "product name already exists")
)

type Service interface {
	ListProducts(ctx context.Context, search string, page, size int) ([]Product, json.Pagination, error)
	CreateProduct(ctx context.Context, name, price string) error
	UpdateProduct(ctx context.Context, id, name, price string) error
	DeleteProduct(ctx context.Context, id string) error
}

type svc struct {
	repo        Repository
	kafkaWriter *kafka.Writer
}

func NewService(repo Repository, kafkaWriter *kafka.Writer) Service {
	return &svc{repo: repo, kafkaWriter: kafkaWriter}
}

func (s *svc) ListProducts(ctx context.Context, search string, page, size int) ([]Product, json.Pagination, error) {
	offset := (page - 1) * size
	items, total, err := s.repo.ListProducts(ctx, search, int32(size), int32(offset))
	if err != nil {
		return nil, json.Pagination{}, err
	}

	return items, json.NewPagination(page, size, total), nil
}

func (s *svc) CreateProduct(ctx context.Context, name, price string) error {
	exists, err := s.repo.CheckNameExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return ErrProductNameExists
	}

	id := ksuid.New().String()

	errCreateProduct := s.repo.CreateProduct(ctx, id, name, price)
	if errCreateProduct != nil {
		return errCreateProduct
	}

	// Emit Kafka Event
	if s.kafkaWriter != nil {
		err := s.kafkaWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(id),
			Value: []byte(name),
		})
		if err != nil {
			fmt.Printf("failed to emit kafka event: %v\n", err)
		}
	}

	return nil
}

func (s *svc) UpdateProduct(ctx context.Context, id, name, price string) error {
	exists, err := s.repo.CheckNameExistsForOther(ctx, name, id)
	if err != nil {
		return err
	}
	if exists {
		return ErrProductNameExists
	}

	return s.repo.UpdateProduct(ctx, id, name, price)
}

func (s *svc) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.DeleteProduct(ctx, id)
}
