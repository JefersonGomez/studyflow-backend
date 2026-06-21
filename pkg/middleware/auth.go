package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"eror": "No tienes autorizacion para esto"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		tokenString, err := ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "no se puedo obtener la informacion (Bearer )"})
			return
		}

		c.Set("userID", tokenString)
		c.Next()

	}
}
