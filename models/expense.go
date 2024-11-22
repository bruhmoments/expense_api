package models

import "time"

// ExpenseRequest represents the structure of the expense creation request body
type ExpenseRequest struct {
	Description string  `json:"description" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
}

// ExpenseResponse represents the structure of the expense data sent to the client
type ExpenseResponse struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}

// ExpenseStats represents the structure of the statistics response
type ExpenseStats struct {
	TotalExpenses float64 `json:"total_expenses"`
	Count         int     `json:"count"`
}
