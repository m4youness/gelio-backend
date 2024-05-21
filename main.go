package main

import (
	"gelio/m/controllers"
	"gelio/m/initializers"
	"gelio/m/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DbConnect()
}

func main() {

	r := gin.Default()

	r.POST("/Login", controllers.Login)
	r.GET("/IsLoggedIn", middleware.RequireAuth, controllers.IsLoggedIn)
	r.POST("SignUp", controllers.SignUp)
	r.Run()
}
