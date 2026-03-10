package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	repo "github.com/sudesh856/ecom-go-api-project/internal/adaptors/postgresql/sqlc"
)


var (
	ErrProductNotFound = errors.New("Product not found.") 
	ErrProductNoStock = errors.New("No stock.")
)

type svc struct {
	repo *repo.Queries
	db *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db: db,

	}
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder CreateOrderParams) (repo.Order, error) {
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer ID is required")
	}

	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("At least one item is needed.")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil{
		return repo.Order{}, err
	}
	for _, item := range tempOrder.Items {
		product, err := qtx.FindProductsByID(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrProductNoStock
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID: order.ID,
			ProductID: item.ProductID,
			Quantity: item.Quantity,
			PriceInRupees: product.PriceInRupees,
		})

		if err != nil {
			return repo.Order{}, err
		}
	}

	tx.Commit(ctx)
		return order, nil
	}	

