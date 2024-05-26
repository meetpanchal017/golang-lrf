package routes

import (
	controller "golang-jwt-demo/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	// incomingRoutes.POST("users/signup", controller.Signup())
	// incomingRoutes.POST("users/login", controller.Login())
	// incomingRoutes.POST("users/send-otp", controller.SendOTP())
	incomingRoutes.POST("login/verify-otp", controller.VerifyOTPNew())
	incomingRoutes.POST("login/phone-number", controller.LoginWithPhoneNumber())
}
