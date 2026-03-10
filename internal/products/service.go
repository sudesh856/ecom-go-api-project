package products

import (
	"context"

	repo "github.com/sudesh856/ecom-go-api-project/internal/adaptors/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProduct(ctx context.Context, id int64) (repo.Product,error)
	CreateProduct(ctx context.Context, arg repo.CreateProductParams) (repo.Product, error)
}
type svc struct {

	repo repo.Querier
}

func NewService(repo repo.Querier) Service{
	return &svc{	
		repo: repo,

	}
}

func(s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func(s *svc) FindProduct(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindProductsByID(ctx, id)
}

func (s *svc) CreateProduct(ctx context.Context, arg repo.CreateProductParams) (repo.Product, error) {
	return s.repo.CreateProduct(ctx, arg)
}

