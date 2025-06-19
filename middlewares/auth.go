package middlewares

import (
	"facegram/database"
	"facegram/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenStr := fields[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			var user models.User
			database.DB.First(&user, claims["user_id"])

			if user.ID == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			c.Set("userID", user.ID)

			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}
