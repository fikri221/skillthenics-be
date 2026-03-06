package orders

import (
	"context"
	"database/sql"
	"errors"
	"nds-go-starter/internal/core/repository"
	"strconv"
)

type Repository interface {
	GetOrderWithItems(ctx context.Context, id string) (Order, error)
	CreateOrder(ctx context.Context, order Order) error
	WithTx(ctx context.Context, fn func(Repository) error) error
}

type repoWrapper struct {
	db    repository.Querier
	rawDB repository.DBTX // To start transactions
}

func NewRepository(db repository.Querier, rawDB repository.DBTX) Repository {
	return &repoWrapper{db: db, rawDB: rawDB}
}

func (r *repoWrapper) WithTx(ctx context.Context, fn func(Repository) error) (err error) {
	db, ok := r.rawDB.(repository.Transactor)
	if !ok {
		return errors.New("repository: database does not support transactions")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var txQuerier repository.Querier = repository.New(tx)
	if _, isLogged := r.rawDB.(*repository.LogDB); isLogged {
		txQuerier = repository.New(repository.NewLogDB(tx))
	}

	txRepo := &repoWrapper{
		db:    txQuerier,
		rawDB: r.rawDB,
	}

	err = fn(txRepo)
	return err
}

func (r *repoWrapper) GetOrderWithItems(ctx context.Context, id string) (Order, error) {
	rows, err := r.db.GetOrderWithItems(ctx, id)
	if err != nil {
		return Order{}, err
	}

	if len(rows) == 0 {
		return Order{}, sql.ErrNoRows
	}

	// Map the flat JOIN result to a nested Order object
	order := Order{
		ID:           rows[0].OrderID,
		CustomerName: rows[0].CustomerName,
		CreatedAt:    rows[0].OrderDate.Time,
	}

	total, _ := strconv.ParseFloat(rows[0].TotalAmount, 64)
	order.TotalAmount = total

	for _, row := range rows {
		if !row.ItemID.Valid {
			continue
		}

		price, _ := strconv.ParseFloat(row.ItemPrice.String, 64)
		// Put items into order, so we can return it as a nested object
		order.Items = append(order.Items, OrderItem{
			ID:          row.ItemID.Int32,
			ProductID:   row.ProductID.String,
			ProductName: row.ProductName.String,
			Quantity:    row.Quantity.Int32,
			Price:       price,
		})
	}

	return order, nil
}

func (r *repoWrapper) CreateOrder(ctx context.Context, order Order) error {
	// This would typically involve a transaction, but for this demo
	// we'll focus on the mapping and JOIN.
	_, err := r.db.CreateOrder(ctx, repository.CreateOrderParams{
		ID:           order.ID,
		CustomerName: order.CustomerName,
		TotalAmount:  strconv.FormatFloat(order.TotalAmount, 'f', 2, 64),
	})
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err := r.db.CreateOrderItem(ctx, repository.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     strconv.FormatFloat(item.Price, 'f', 2, 64),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
