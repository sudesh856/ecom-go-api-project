-- name: ListProducts :many
SELECT * FROM products;

-- name: FindProductsByID :one
SELECT * FROM products WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (
    customer_id
  )  VALUES ($1) RETURNING *;


-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_in_rupees) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateProduct :one
INSERT INTO products (name,  price_in_rupees, quantity)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1; 