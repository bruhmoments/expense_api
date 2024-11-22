package middleware

import (
	"net/http"
	"os"
	"pengeluaran/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims = models.JWTClaims

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*JWTClaims); ok {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
		}
		// fmt.Println("Gin Middleware context:", c)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	allowedURL := os.Getenv("FRONTEND_URL")
	if allowedURL == "" {
		allowedURL = "*"
	}

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedURL)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // Respond with 204 for preflight request
			return
		}

		c.Next()
	}
}
