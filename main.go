package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"pengeluaran/config"
	"pengeluaran/db/queries"
	"pengeluaran/handlers"
	"pengeluaran/middleware"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in the environment file")
	}

	fmt.Println("JWT_SECRET loaded:", jwtSecret)

	// Load config
	dbConfig := config.LoadConfig()

	// Debug print to console
	fmt.Printf("Loaded Config: %+v\n", dbConfig) // This will print the entire struct

	// Build the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize queries Gin router
	q := queries.New(db)
	r := gin.Default()

	// Apply CORS middleware globally
	r.Use(middleware.CORSMiddleware())

	// Authentication routes (no middleware)
	auth := r.Group("/")
	{
		auth.OPTIONS("/login", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})
		auth.OPTIONS("/register", func(c *gin.Context) {
			c.Status(http.StatusNoContent) // Handle OPTIONS preflight request
		})
		auth.POST("/register", handlers.RegisterHandler(q))
		auth.POST("/login", handlers.LoginHandler(q))
	}

	// Protected routes for expenses (middleware applied)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret)) // Apply middleware here
	{
		protected.GET("/expenses", handlers.GetExpensesHandler(q))
		protected.POST("/expenses", handlers.CreateExpenseHandler(q))
		protected.PUT("/expenses", handlers.UpdateExpenseHandler(q))
		protected.DELETE("/expenses", handlers.DeleteExpenseHandler(q))
		protected.GET("/expenses/stats", handlers.GetExpenseStatsHandler(q))
	}

	// Start the server
	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")

}
