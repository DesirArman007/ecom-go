package products

import (
	"context"
	"errors"
	"strings"

	repo "github.com/DesirArman007/ecom/internal/adapters/postgresql/sqlc"
)

var (
	ErrProductNameRequired = errors.New("product name is required")
	ErrPriceInvalid        = errors.New("price cannot be negative")
	ErrQuantityInvalid     = errors.New("quantity cannot be negative")
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductByID(ctx context.Context, id int64) (repo.Product, error)
	CreateProduct(ctx context.Context, params createProductParams) (repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) FindProductByID(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindProductByID(ctx, id)
}

func (s *svc) CreateProduct(ctx context.Context, params createProductParams) (repo.Product, error) {

	// validate payload
	if strings.TrimSpace(params.Name) == "" {
		return repo.Product{}, ErrProductNameRequired
	}

	if params.PriceInCents < 0 {
		return repo.Product{}, ErrPriceInvalid
	}

	if params.Quantity < 0 {
		return repo.Product{}, ErrQuantityInvalid
	}

	return s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:         params.Name,
		PriceInCents: params.PriceInCents,
		Quantity:     params.Quantity,
	})
}
