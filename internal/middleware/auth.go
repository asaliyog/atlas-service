package middleware

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang-service/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AzureJWKS represents Azure's JSON Web Key Set
type AzureJWKS struct {
	Keys []AzureJWK `json:"keys"`
}

// AzureJWK represents a single JSON Web Key
type AzureJWK struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
	X5t string   `json:"x5t"`
}

// AzureTokenClaims represents Azure token claims
type AzureTokenClaims struct {
	Aud           string `json:"aud"`
	Iss           string `json:"iss"`
	Iat           int64  `json:"iat"`
	Nbf           int64  `json:"nbf"`
	Exp           int64  `json:"exp"`
	Aio           string `json:"aio"`
	Appid         string `json:"appid"`
	Appidacr      string `json:"appidacr"`
	Idp           string `json:"idp"`
	ObjID         string `json:"oid"`
	Rh            string `json:"rh"`
	Sub           string `json:"sub"`
	TenantID      string `json:"tid"`
	UTI           string `json:"uti"`
	Ver           string `json:"ver"`
	Roles         []string `json:"roles"`
	jwt.RegisteredClaims
}

var (
	// Cache for Azure public keys
	azurePublicKeys map[string]*rsa.PublicKey
	keysLastFetch   time.Time
	keysCacheTTL    = 24 * time.Hour
)

// AzureEntraAuth validates Azure Entra ID tokens with client credential flow
func AzureEntraAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bypass authentication for local development
		if cfg.BypassAuth {
			log.Println("Auth bypassed for local development")
			c.Set("user_id", "dev-user")
			c.Set("tenant_id", "dev-tenant")
			c.Set("client_id", "dev-client")
			c.Next()
			return
		}

		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format. Expected 'Bearer <token>'"})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &AzureTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return getAzurePublicKey(token, cfg)
		})

		if err != nil {
			log.Printf("Token validation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract and validate claims
		if claims, ok := token.Claims.(*AzureTokenClaims); ok && token.Valid {
			// Validate issuer
			expectedIssuer := fmt.Sprintf("https://sts.windows.net/%s/", cfg.AzureTenantID)
			if claims.Iss != expectedIssuer {
				log.Printf("Invalid issuer: expected %s, got %s", expectedIssuer, claims.Iss)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token issuer"})
				c.Abort()
				return
			}

			// Validate audience (should be the client ID)
			if claims.Aud != cfg.AzureClientID {
				log.Printf("Invalid audience: expected %s, got %s", cfg.AzureClientID, claims.Aud)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token audience"})
				c.Abort()
				return
			}

			// Validate tenant ID
			if claims.TenantID != cfg.AzureTenantID {
				log.Printf("Invalid tenant ID: expected %s, got %s", cfg.AzureTenantID, claims.TenantID)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid tenant ID"})
				c.Abort()
				return
			}

			// Store token information in context
			c.Set("user_id", claims.Sub)
			c.Set("tenant_id", claims.TenantID)
			c.Set("client_id", claims.Appid)
			c.Set("roles", claims.Roles)
			c.Set("object_id", claims.ObjID)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getAzurePublicKey retrieves the public key for token validation
func getAzurePublicKey(token *jwt.Token, cfg *config.Config) (interface{}, error) {
	// Check if we need to refresh the keys
	if azurePublicKeys == nil || time.Since(keysLastFetch) > keysCacheTTL {
		if err := refreshAzurePublicKeys(cfg); err != nil {
			return nil, fmt.Errorf("failed to refresh Azure public keys: %w", err)
		}
	}

	// Get the key ID from the token header
	kidInterface, ok := token.Header["kid"]
	if !ok {
		return nil, fmt.Errorf("token header missing 'kid' claim")
	}

	kid, ok := kidInterface.(string)
	if !ok {
		return nil, fmt.Errorf("'kid' claim is not a string")
	}

	// Look up the public key
	publicKey, exists := azurePublicKeys[kid]
	if !exists {
		// Try to refresh keys once more
		if err := refreshAzurePublicKeys(cfg); err != nil {
			return nil, fmt.Errorf("failed to refresh Azure public keys: %w", err)
		}
		
		publicKey, exists = azurePublicKeys[kid]
		if !exists {
			return nil, fmt.Errorf("unable to find public key for kid: %s", kid)
		}
	}

	return publicKey, nil
}

// refreshAzurePublicKeys fetches the latest public keys from Azure
func refreshAzurePublicKeys(cfg *config.Config) error {
	jwksURL := fmt.Sprintf("https://login.microsoftonline.com/%s/discovery/v2.0/keys", cfg.AzureTenantID)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", jwksURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch JWKS, status: %d", resp.StatusCode)
	}

	var jwks AzureJWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	// Convert JWK to RSA public keys
	publicKeys := make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		if key.Kty == "RSA" && key.Use == "sig" {
			publicKey, err := jwkToRSAPublicKey(key)
			if err != nil {
				log.Printf("Failed to convert JWK to RSA public key for kid %s: %v", key.Kid, err)
				continue
			}
			publicKeys[key.Kid] = publicKey
		}
	}

	azurePublicKeys = publicKeys
	keysLastFetch = time.Now()
	
	log.Printf("Successfully loaded %d Azure public keys", len(publicKeys))
	return nil
}

// jwkToRSAPublicKey converts a JWK to an RSA public key
func jwkToRSAPublicKey(jwk AzureJWK) (*rsa.PublicKey, error) {
	if len(jwk.X5c) == 0 {
		return nil, fmt.Errorf("no x5c certificate found in JWK")
	}

	// Decode the x5c certificate
	certData, err := base64.StdEncoding.DecodeString(jwk.X5c[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode x5c certificate: %w", err)
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Extract the public key
	publicKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("certificate does not contain an RSA public key")
	}

	return publicKey, nil
}