// models/jwt_claims.go
package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
