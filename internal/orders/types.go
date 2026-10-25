package orders

import repo "github.com/DesirArman007/ecom/internal/adapters/postgresql/sqlc"

type orderItem struct {
	ProductID int32 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type createOrderParams struct {
	CustomerId int64       `json:"customerId"`
	Items      []orderItem `json:"items"`
}

type OrderResponse struct {
	Order           repo.Order                       `json:"order"`
	Items           []repo.GetOrderItemsByOrderIDRow `json:"items"`
	TotalPriceCents int64                            `json:"total_price_cents"`
}
