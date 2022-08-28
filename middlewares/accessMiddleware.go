package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func HasStaffAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		if Contains(c.GetStringSlice("userRole"), "student") {
			c.JSON(http.StatusForbidden, gin.H{"status": "failed", "error": "Forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func HasAdminAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !(Contains(c.GetStringSlice("userRole"), "admin")) {
			c.JSON(http.StatusForbidden, gin.H{"status": "failed", "error": "Forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
