// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: expenses.sql

package queries

import (
	"context"
	"database/sql"
	"time"
)

const createExpense = `-- name: CreateExpense :one
INSERT INTO expenses (user_id, description, amount) 
VALUES ($1, $2, $3)
RETURNING id, description, amount, created_at
`

type CreateExpenseParams struct {
	UserID      int32
	Description string
	Amount      string
}

type CreateExpenseRow struct {
	ID          int32
	Description string
	Amount      string
	CreatedAt   sql.NullTime
}

func (q *Queries) CreateExpense(ctx context.Context, arg CreateExpenseParams) (CreateExpenseRow, error) {
	row := q.db.QueryRowContext(ctx, createExpense, arg.UserID, arg.Description, arg.Amount)
	var i CreateExpenseRow
	err := row.Scan(
		&i.ID,
		&i.Description,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteExpense = `-- name: DeleteExpense :exec
DELETE FROM expenses 
WHERE id = $1 AND user_id = $2
`

type DeleteExpenseParams struct {
	ID     int32
	UserID int32
}

func (q *Queries) DeleteExpense(ctx context.Context, arg DeleteExpenseParams) error {
	_, err := q.db.ExecContext(ctx, deleteExpense, arg.ID, arg.UserID)
	return err
}

const getAllExpensesByUser = `-- name: GetAllExpensesByUser :many
SELECT id, user_id, description, amount, created_at 
FROM expenses 
WHERE user_id = $1
  AND created_at >= $2 
  AND created_at < $3
ORDER BY created_at DESC
`

type GetAllExpensesByUserParams struct {
	UserID      int32
	CreatedAt   time.Time
	CreatedAt_2 time.Time
}

func (q *Queries) GetAllExpensesByUser(ctx context.Context, arg GetAllExpensesByUserParams) ([]Expense, error) {
	rows, err := q.db.QueryContext(ctx, getAllExpensesByUser, arg.UserID, arg.CreatedAt, arg.CreatedAt_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Expense
	for rows.Next() {
		var i Expense
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Description,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getExpenseStatistics = `-- name: GetExpenseStatistics :many
SELECT 
    DATE_TRUNC('month', created_at) AS month,
    SUM(amount) AS total_amount
FROM expenses
WHERE user_id = $1 
  AND created_at >= $2 
  AND created_at < $3
GROUP BY month
ORDER BY month ASC
`

type GetExpenseStatisticsParams struct {
	UserID      int32
	CreatedAt   time.Time
	CreatedAt_2 time.Time
}

type GetExpenseStatisticsRow struct {
	Month       time.Time
	TotalAmount int64
}

func (q *Queries) GetExpenseStatistics(ctx context.Context, arg GetExpenseStatisticsParams) ([]GetExpenseStatisticsRow, error) {
	rows, err := q.db.QueryContext(ctx, getExpenseStatistics, arg.UserID, arg.CreatedAt, arg.CreatedAt_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetExpenseStatisticsRow
	for rows.Next() {
		var i GetExpenseStatisticsRow
		if err := rows.Scan(&i.Month, &i.TotalAmount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateExpense = `-- name: UpdateExpense :one
UPDATE expenses
SET description = $2, amount = $3
WHERE id = $1 AND user_id = $4
RETURNING id, description, amount, created_at
`

type UpdateExpenseParams struct {
	ID          int32
	Description string
	Amount      string
	UserID      int32
}

type UpdateExpenseRow struct {
	ID          int32
	Description string
	Amount      string
	CreatedAt   sql.NullTime
}

func (q *Queries) UpdateExpense(ctx context.Context, arg UpdateExpenseParams) (UpdateExpenseRow, error) {
	row := q.db.QueryRowContext(ctx, updateExpense,
		arg.ID,
		arg.Description,
		arg.Amount,
		arg.UserID,
	)
	var i UpdateExpenseRow
	err := row.Scan(
		&i.ID,
		&i.Description,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
