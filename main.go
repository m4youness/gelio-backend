package main

import (
	"gelio/m/controllers"
	"gelio/m/initializers"
	"gelio/m/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DbConnect()
}

func main() {

	r := gin.Default()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	})

	r.Use(corsHandler(c))
	r.POST("/SignUp", controllers.SignUp)
	r.POST("/SignIn", controllers.Login)
	r.GET("/IsAuthenticated", middleware.RequireAuth, controllers.IsLoggedIn)
	r.GET("/Logout", middleware.RequireAuth, controllers.Logout)
	r.Run()
}

func corsHandler(corsMiddleware *cors.Cors) gin.HandlerFunc {
	return func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	}
}
