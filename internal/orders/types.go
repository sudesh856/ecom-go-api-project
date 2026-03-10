package orders

import (
	"context"

	repo "github.com/sudesh856/ecom-go-api-project/internal/adaptors/postgresql/sqlc"
)

type orderItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type CreateOrderParams struct {
	CustomerID int64       `json:"customerId"`
	Items      []orderItem `json:"items"`
}

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder CreateOrderParams) (repo.Order, error)
}