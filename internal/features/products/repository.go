package products

import (
	"context"
	"database/sql"
	"nds-go-starter/internal/core/repository"
)

type Repository interface {
	ListProducts(ctx context.Context, search string, limit, offset int32) ([]Product, int64, error)
	GetProductByID(ctx context.Context, id string) (Product, error)
	CreateProduct(ctx context.Context, id, name, price string) error
	UpdateProduct(ctx context.Context, id, name, price string) error
	DeleteProduct(ctx context.Context, id string) error
	CheckNameExists(ctx context.Context, name string) (bool, error)
	CheckNameExistsForOther(ctx context.Context, name, id string) (bool, error)
}

type repoWrapper struct {
	db repository.Querier
}

func NewRepository(db repository.Querier) Repository {
	return &repoWrapper{db: db}
}

func (r *repoWrapper) ListProducts(ctx context.Context, search string, limit, offset int32) ([]Product, int64, error) {
	total, err := r.db.CountProducts(ctx, repository.CountProductsParams{
		Search: search,
	})
	if err != nil {
		return nil, 0, err
	}

	items, err := r.db.ListProducts(ctx, repository.ListProductsParams{
		Search: search,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, err
	}

	var products []Product
	for _, item := range items {
		products = append(products, mapToDomain(item))
	}

	return products, total, nil
}

func (r *repoWrapper) GetProductByID(ctx context.Context, id string) (Product, error) {
	item, err := r.db.GetProductByID(ctx, sql.NullString{String: id, Valid: true})
	if err != nil {
		return Product{}, err
	}
	return mapToDomain(item), nil
}

func mapToDomain(m repository.MsProduct) Product {
	return Product{
		ID:        m.ID.String,
		Name:      m.Name,
		Price:     m.Price,
		RecStatus: m.RecStatus,
		CreatedAt: m.CreatedAt.Time,
		UpdatedAt: m.UpdatedAt.Time,
	}
}

func (r *repoWrapper) CreateProduct(ctx context.Context, id, name, price string) error {
	_, err := r.db.CreateProduct(ctx, repository.CreateProductParams{
		ID:    sql.NullString{String: id, Valid: true},
		Name:  name,
		Price: price,
	})
	return err
}

func (r *repoWrapper) UpdateProduct(ctx context.Context, id, name, price string) error {
	_, err := r.db.UpdateProduct(ctx, repository.UpdateProductParams{
		ID:    sql.NullString{String: id, Valid: true},
		Name:  name,
		Price: price,
	})
	return err
}

func (r *repoWrapper) DeleteProduct(ctx context.Context, id string) error {
	_, err := r.db.DeleteProduct(ctx, sql.NullString{String: id, Valid: true})
	return err
}

func (r *repoWrapper) CheckNameExists(ctx context.Context, name string) (bool, error) {
	return r.db.CheckProductNameExists(ctx, name)
}

func (r *repoWrapper) CheckNameExistsForOther(ctx context.Context, name, id string) (bool, error) {
	return r.db.CheckProductNameExistsForOther(ctx, repository.CheckProductNameExistsForOtherParams{
		Name: name,
		ID:   sql.NullString{String: id, Valid: true},
	})
}
