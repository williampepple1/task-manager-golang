package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT key used to create the signature.
// This should ideally come from a more secure source such as environment variables.
var jwtKey = []byte("secret_key")

// Authorize ensures the token is valid and sets the user's information in the context.
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		tokenStr := c.GetHeader("Authorization")

		claims := &jwt.StandardClaims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort() // abort further handlers
			return
		}

		// Set the username in the context for subsequent handlers to use
		c.Set("username", claims.Subject)
		c.Next() // proceed to next middleware or handler
	}
}
