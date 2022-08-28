package middlewares

import (
	"fmt"
	helper "login/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken, _ := c.Cookie("token")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": "No Auth Header provided"})
			c.Abort()
			return
		}
		claims, errr := helper.TokenHelper.ValidateToken(clientToken)
		if errr != "" {
			fmt.Println("validateTokenError")
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errr})
			c.Abort()
			return
		}
		c.Set("name", claims.Name)
		c.Set("userId", claims.UserId)
		c.Set("userRole", claims.UserRole)
		c.Next()

	}
}
