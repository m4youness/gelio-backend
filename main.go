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
	r.POST("/User", controllers.SignUp)
	r.POST("/SignIn", controllers.Login)
	r.GET("/IsAuthenticated", middleware.RequireAuth, controllers.IsLoggedIn)
	r.GET("/Logout", middleware.RequireAuth, controllers.Logout)
	r.POST("/Person", controllers.AddPerson)
	r.GET("/UserId", controllers.GetUserId)
	r.GET("/User/:id", controllers.GetUser)
	r.POST("/User/Exists", controllers.DoesUserExist)
	r.GET("/Person/:id", controllers.GetPerson)
	r.GET("/Countries", controllers.GetAllCountries)
	r.POST("/GetCountryWithName", controllers.GetCountryIdWithName)
	r.GET("/Country/:id", controllers.GetCountryNameWithId)
	r.POST("/Image", controllers.AddImage)
	r.GET("/LoadMessages/:id", controllers.LoadMessages)
	r.Run()

}

func corsHandler(corsMiddleware *cors.Cors) gin.HandlerFunc {
	return func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	}
}
