-- name: ListProducts :many
SELECT * FROM products;

-- name: FindProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (
    customer_id
) VALUES ($1) RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_cents) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetOrderItemsByOrderID :many
SELECT
    oi.id,
    oi.order_id,
    oi.product_id,
    p.name AS product_name,
    oi.quantity,
    oi.price_cents,
    oi.created_at
FROM order_items oi
JOIN products p ON p.id = oi.product_id
WHERE oi.order_id = $1
ORDER BY oi.id;

-- name: DecreaseProductStock :one
UPDATE products
SET 
    quantity = quantity - $2,
    updated_at = NOW()
Where id = $1 AND quantity >= $2
RETURNING *;

-- name: CreateProduct :one
INSERT INTO products (
    name,
    price_in_cents,
    quantity
)
VALUES ($1, $2, $3)
RETURNING *;

