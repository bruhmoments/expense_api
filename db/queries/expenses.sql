-- name: CreateExpense :one
INSERT INTO expenses (user_id, description, amount) 
VALUES ($1, $2, $3)
RETURNING id, description, amount, created_at;

-- name: GetAllExpensesByUser :many
SELECT id, user_id, description, amount, created_at 
FROM expenses 
WHERE user_id = $1
  AND created_at >= $2 
  AND created_at < $3;

-- name: UpdateExpense :one
UPDATE expenses
SET description = $2, amount = $3
WHERE id = $1 AND user_id = $4
RETURNING id, description, amount, created_at;

-- name: DeleteExpense :exec
DELETE FROM expenses 
WHERE id = $1 AND user_id = $2;

-- name: GetExpenseStatistics :many
SELECT 
    DATE_TRUNC('month', created_at) AS month,
    SUM(amount) AS total_amount
FROM expenses
WHERE user_id = $1 
  AND created_at >= $2 
  AND created_at < $3
GROUP BY month
ORDER BY month DESC;
