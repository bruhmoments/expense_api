package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"pengeluaran/db/queries" // Update this to match your module name

	"github.com/gin-gonic/gin"
)

type CreateExpenseRequest struct {
	Description string `json:"description" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
}

func CreateExpenseHandler(q *queries.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateExpenseRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		// fmt.Println("Gin Middleware context:", c)
		// Extract user ID from JWT token (to be added in the auth middleware)
		userID, _ := GetUserID(c)

		expense, err := q.CreateExpense(c, queries.CreateExpenseParams{
			UserID:      userID,
			Description: req.Description,
			Amount:      req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Data created successfully",
			"data": gin.H{
				"id":          expense.ID,
				"description": expense.Description,
				"amount":      expense.Amount,
				"created_at":  expense.CreatedAt,
			},
		})
	}
}

type GetExpensesResponse struct {
	ID          int32     `json:"id"`
	Description string    `json:"description"`
	Amount      string    `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}

func GetExpensesHandler(q *queries.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user ID from JWT token (to be added in the auth middleware)
		userID, _ := GetUserID(c)

		startDateParam := c.Query("start_date")
		endDateParam := c.Query("end_date")

		defaultStartDate := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
		defaultEndDate := time.Date(time.Now().Year(), time.Now().Month(), 31, 23, 59, 59, 0, time.UTC)

		var startDate time.Time
		if startDateParam == "" {
			startDate = defaultStartDate
		} else {
			startDate, _ = time.Parse("2006-01-02", startDateParam)
		}

		var endDate time.Time
		if endDateParam == "" {
			endDate = defaultEndDate
		} else {
			parsedEndDate, _ := time.Parse("2006-01-02", endDateParam)
			endDate = time.Date(
				parsedEndDate.Year(),
				parsedEndDate.Month(),
				parsedEndDate.Day(),
				23, 59, 59, 999999999,
				parsedEndDate.Location(),
			)
		}

		expenses, err := q.GetAllExpensesByUser(c, queries.GetAllExpensesByUserParams{
			UserID:      userID,
			CreatedAt:   startDate,
			CreatedAt_2: endDate,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
			return
		}

		if len(expenses) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No data found", "data": []string{}})
			return
		}

		var response []GetExpensesResponse
		for _, expense := range expenses {
			response = append(response, GetExpensesResponse{
				ID:          expense.ID,
				Description: expense.Description,
				Amount:      expense.Amount,
				CreatedAt:   expense.CreatedAt.Time,
			})
		}

		c.JSON(http.StatusOK, gin.H{"message": "Data retrieved successfully", "data": response})
	}
}

type UpdateExpenseRequest struct {
	ID          int32  `json:"id" binding:"required"`
	Description string `json:"description" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
}

func UpdateExpenseHandler(q *queries.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateExpenseRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Extract user ID from JWT token (to be added in the auth middleware)
		userID, _ := GetUserID(c)

		// Update expense in the database
		expense, err := q.UpdateExpense(c, queries.UpdateExpenseParams{
			ID:          req.ID,
			UserID:      userID,
			Description: req.Description,
			Amount:      req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Data updated successfully",
			"data": gin.H{
				"id":          expense.ID,
				"description": expense.Description,
				"amount":      expense.Amount,
				"created_at":  expense.CreatedAt,
			},
		})
	}
}

type DeleteExpenseRequest struct {
	ID int32 `json:"id" binding:"required"`
}

func DeleteExpenseHandler(q *queries.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DeleteExpenseRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Extract user ID from JWT token (to be added in the auth middleware)
		userID, _ := GetUserID(c)

		// Delete expense from the database
		err := q.DeleteExpense(c, queries.DeleteExpenseParams{
			ID:     req.ID,
			UserID: userID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Expense deleted"})
	}
}

type GetExpenseStatsResponse struct {
	Month       string `json:"month"`
	TotalAmount string `json:"total_amount"`
}

func GetExpenseStatsHandler(q *queries.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user ID from JWT token (to be added in the auth middleware)
		userID, _ := GetUserID(c)

		startDateParam := c.Query("start_date")
		endDateParam := c.Query("end_date")

		currentYear := time.Now().Year()
		defaultStartDate := time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.UTC)
		defaultEndDate := time.Date(currentYear, time.December, 31, 23, 59, 59, 0, time.UTC)

		var startDate time.Time
		if startDateParam == "" {
			startDate = defaultStartDate
		} else {
			startDate, _ = time.Parse("2006-01-02", startDateParam)
		}

		var endDate time.Time
		if endDateParam == "" {
			endDate = defaultEndDate
		} else {
			parsedEndDate, _ := time.Parse("2006-01-02", endDateParam)
			endDate = time.Date(
				parsedEndDate.Year(),
				parsedEndDate.Month(),
				parsedEndDate.Day(),
				23, 59, 59, 999999999, // End of the day
				parsedEndDate.Location(),
			)
		}

		stats, err := q.GetExpenseStatistics(c, queries.GetExpenseStatisticsParams{
			UserID:      userID,
			CreatedAt:   startDate,
			CreatedAt_2: endDate,
		})

		if err != nil {
			fmt.Println("error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expense statistics"})
			return
		}

		if len(stats) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No data found", "data": []string{}})
			return
		}

		var response []GetExpenseStatsResponse
		for _, stat := range stats {
			response = append(response, GetExpenseStatsResponse{
				Month:       stat.Month.Format("01-2006"),
				TotalAmount: fmt.Sprintf("%d", stat.TotalAmount),
			})
		}

		// c.JSON(http.StatusOK, response)
		c.JSON(http.StatusOK, gin.H{"message": "Data retrieved successfully", "data": response})
	}
}

func GetUserID(c *gin.Context) (int32, error) {
	// Retrieve user_id from the context
	userID, _ := c.Get("user_id")

	switch v := userID.(type) {
	case int:
		return int32(v), nil
	case int32:
		return v, nil
	case string:
		parsedID, _ := strconv.ParseInt(v, 10, 32)
		return int32(parsedID), nil
	default:
		return 0, nil
	}
}
