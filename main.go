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
	// User
	r.POST("/User", controllers.SignUp)
	r.POST("/SignIn", controllers.Login)
	r.GET("/IsAuthenticated", middleware.RequireAuth, controllers.IsLoggedIn)
	r.GET("/Logout", middleware.RequireAuth, controllers.Logout)
	r.GET("/User/:id", controllers.GetUser)
	r.POST("/User/Exists", controllers.DoesUserExist)
	r.GET("/UserId", controllers.GetUserId)

	// People
	r.POST("/Person", controllers.AddPerson)
	r.GET("/Person/:id", controllers.GetPerson)

	// Country
	r.GET("/Countries", controllers.GetAllCountries)
	r.POST("/GetCountryWithName", controllers.GetCountryIdWithName)
	r.GET("/Country/:id", controllers.GetCountryNameWithId)

	// Image
	r.POST("/Image", controllers.AddImage)

	// Message
	r.GET("/LoadContacts/:id", controllers.LoadContacts)
	r.POST("/LoadMessages", controllers.LoadMessages)
	r.GET("/MessageInfo/:id", controllers.GetMessageInfoFromId)
	r.POST("/Message", controllers.SendMessage)
	r.POST("/Contact", controllers.AddContact)

	r.Run()

}

func corsHandler(corsMiddleware *cors.Cors) gin.HandlerFunc {
	return func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	}
}
