package routes

import (
	"login/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.PUT("/api/users/forgotPassword", controllers.ResetPassword)
	incomingRoutes.POST("/api/students/signup", controllers.RegisterStudent)
	incomingRoutes.POST("/api/users/login", controllers.Login)
	incomingRoutes.POST("/api/users/logout", controllers.Logout)

}
