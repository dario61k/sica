package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sica/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")

		check := strings.Split(token, " ")
		if len(check) != 2 || check[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Credentials not provided",
			})
			c.Abort()
			return
		}

		claim, err := jwt.ValidateToken(check[1])

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Credentials not provided",
			})
			c.Abort()
			return
		}

		uuid_token, _ := uuid.Parse(claim.Subject)

		c.Set("user_id", uuid_token)

		c.Next()
	}
}
