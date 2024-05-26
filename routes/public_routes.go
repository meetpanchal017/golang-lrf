package routes

import (
	controller "golang-jwt-demo/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/uploads/:id", controller.GetImage())
	incomingRoutes.GET("/favicon.ico", func(c *gin.Context) {

		c.Status(http.StatusNoContent)
	})
}
