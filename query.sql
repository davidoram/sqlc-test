
-- name: CreateCustomer :one
INSERT INTO customers (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: UpdateCustomer :one
UPDATE customers
SET name = $2
WHERE id = $1
RETURNING *;

-- name: GetCustomerByID :one
SELECT * FROM customers
WHERE id = $1;


-- name: GetCustomerRevisions :many
SELECT * FROM customers
WHERE id = $1
UNION ALL
SELECT * FROM customer_revisions
WHERE customer_id = $1
ORDER BY revision DESC;

