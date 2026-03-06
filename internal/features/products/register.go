package products

import (
	"nds-go-starter/internal/core/repository"

	"github.com/go-chi/chi/v5"
	"github.com/segmentio/kafka-go"
)

func Register(r chi.Router, db repository.DBTX, kafkaWriter *kafka.Writer) {
	querier := repository.New(db)
	repo := NewRepository(querier)
	service := NewService(repo, kafkaWriter)
	handler := NewHandler(service)

	r.Route("/products", func(r chi.Router) {
		r.Get("/", handler.ListProducts)
		r.Post("/", handler.CreateProduct)
		r.Put("/{id}", handler.UpdateProduct)
		r.Delete("/{id}", handler.DeleteProduct)
	})
}
