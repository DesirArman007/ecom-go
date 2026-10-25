package orders

import (
	"context"
	"errors"
	"fmt"

	repo "github.com/DesirArman007/ecom/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound = errors.New("Product not found")
	ErrProductNoStock  = errors.New("Product has not enough stock")
	ErrOrderNotFound   = errors.New("Order not found")
)

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
	GetOrderByID(ctx context.Context, id int64) (OrderResponse, error)
}

// TxQuerier can do everything Querier can do, and it can also use transactions
type TxQuerier interface {
	repo.Querier
	WithTx(tx pgx.Tx) *repo.Queries
}

type svc struct {
	repo TxQuerier
	db   *pgx.Conn
}

func NewService(repo TxQuerier, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) GetOrderByID(ctx context.Context, id int64) (OrderResponse, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return OrderResponse{}, ErrOrderNotFound
		}
		return OrderResponse{}, err
	}

	items, err := s.repo.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return OrderResponse{}, err
	}

	var totalPriceCents int64
	for _, item := range items {
		totalPriceCents += int64(item.Quantity) * int64(item.PriceCents)
	}

	return OrderResponse{
		Order:           order,
		Items:           items,
		TotalPriceCents: totalPriceCents,
	}, nil
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	//  validation payload
	if tempOrder.CustomerId == 0 {
		return repo.Order{}, fmt.Errorf("Customer ID is required")
	}

	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("At least one item is required")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	// If error occures rollback
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// 1. Create a new order in the database
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerId)

	if err != nil {
		return repo.Order{}, err
	}

	// 2. look for the product if exist
	for _, item := range tempOrder.Items {
		product, err := qtx.FindProductByID(ctx, int64(item.ProductID))
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrProductNoStock
		}

		// 3. Create order
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  int64(item.ProductID),
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCents,
		})

		if err != nil {
			return repo.Order{}, err
		}

		// Update the product stock quantity
		_, err = qtx.DecreaseProductStock(ctx, repo.DecreaseProductStockParams{
			ID:       int64(item.ProductID),
			Quantity: item.Quantity,
		})

		if err != nil {
			return repo.Order{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, err
	}

	return order, nil
}
