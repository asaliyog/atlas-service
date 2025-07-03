package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"golang-service/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AzureEntraAuth validates Azure Entra ID tokens
func AzureEntraAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Store user information in context
			c.Set("user_id", claims["sub"])
			c.Set("user_email", claims["email"])
			c.Set("tenant_id", claims["tid"])
		}

		c.Next()
	}
}

// Note: In a real implementation, you would:
// 1. Validate the token against Azure Entra ID's public keys
// 2. Verify the audience (aud) claim matches your application
// 3. Check the issuer (iss) claim
// 4. Validate expiration time
// This simplified version is for demonstration purposes