package routes

import (
	"login/controllers"
	middleware "login/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/api/activeUser", controllers.GetActiveUser)
	incomingRoutes.GET("/api/users/:id", controllers.GetUser)
	incomingRoutes.GET("/api/roles", controllers.GetAllRoles)
	incomingRoutes.PUT("/api/users/:id", controllers.UpdateUser)
	incomingRoutes.PUT("/api/users/changePassword/:id", controllers.ChangePassword)

	incomingRoutes.Use(middleware.HasStaffAccess())
	incomingRoutes.GET("/api/students", controllers.GetAllStudents)

	incomingRoutes.Use(middleware.HasAdminAccess())

	incomingRoutes.GET("/api/users", controllers.GetAllUser)
	incomingRoutes.POST("/api/users/register", controllers.CreateUser)
	incomingRoutes.DELETE("/api/users/:id", controllers.DeleteUser)
	incomingRoutes.POST("/api/roles", controllers.CreateRole)
	incomingRoutes.PUT("/api/roles/:id", controllers.UpdateRole)
	incomingRoutes.DELETE("/api/roles/:id", controllers.DeleteRole)
	incomingRoutes.GET("/api/roles/:id", controllers.GetRole)
}
