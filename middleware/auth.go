package middleware

import (
	"strings"

	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

// JWT key used to create the signature.
// var jwtKey = os.Getenv("JWT_SECRET_KEY")
var jwtKey = "JWT_SECRET_KEY"

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header not provided"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		claims := &jwt.StandardClaims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set the username in the context for subsequent handlers to use
		c.Set("username", claims.Subject)
		c.Next() // proceed to next middleware or handler
	}
}
