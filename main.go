package main

import (
	"gelio/m/controllers"
	"gelio/m/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DbConnect()
}

func main() {

	r := gin.Default()

	r.POST("/Login", controllers.Login)

	r.Run()
}
