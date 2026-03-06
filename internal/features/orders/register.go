package orders

import (
	"nds-go-starter/internal/core/repository"

	"github.com/go-chi/chi/v5"
	"github.com/segmentio/kafka-go"
)

func Register(r chi.Router, db repository.DBTX, kafkaWriter *kafka.Writer) {
	querier := repository.New(db)
	repo := NewRepository(querier, db)
	service := NewService(repo, kafkaWriter)
	handler := NewHandler(service)

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", handler.CreateOrder)
		r.Get("/{id}", handler.GetOrder)
	})
}
