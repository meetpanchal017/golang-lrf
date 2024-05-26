package routes

import (
	controller "golang-jwt-demo/controllers"
	"golang-jwt-demo/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:uid", controller.GetUser())
	incomingRoutes.PUT("users/update", controller.UpdateUser())
	incomingRoutes.POST("users/update-profile-image", controller.UpdateProfileImage())
}
