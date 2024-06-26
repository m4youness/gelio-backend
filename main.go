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
	initializers.CloudinaryConnect()
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
	r.GET("/Logout", middleware.RequireAuth, controllers.Logout)
	r.GET("/User/:id", middleware.RequireAuth, controllers.GetUser)
	r.POST("/User/Exists", controllers.DoesUserExist)
	r.GET("/UserId", middleware.RequireAuth, controllers.GetUserId)
	r.GET("/User/InActive/:id", middleware.RequireAuth, controllers.MakeUserInActive)
	r.GET("/User/IsNotActive/:username", controllers.UserActivity)

	// People
	r.POST("/Person", controllers.AddPerson)
	r.GET("/Person/:id", middleware.RequireAuth, controllers.GetPerson)

	// Country
	r.GET("/Countries", controllers.GetAllCountries)
	r.POST("/GetCountryWithName", controllers.GetCountryIdWithName)
	r.GET("/Country/:id", middleware.RequireAuth, controllers.GetCountryNameWithId)

	// Image
	r.POST("/Image", controllers.UploadImage)
	r.GET("/Image/:id", middleware.RequireAuth, controllers.FindImage)

	// Message
	r.GET("/LoadContacts/:id", middleware.RequireAuth, controllers.LoadContacts)
	r.POST("/LoadMessages", middleware.RequireAuth, controllers.LoadMessages)
	r.POST("/Message", middleware.RequireAuth, controllers.SendMessage)
	r.POST("/Contact", middleware.RequireAuth, controllers.AddContact)
	r.POST("/ContactExist", middleware.RequireAuth, controllers.IsPersonNotContact)

	r.Run()

}

func corsHandler(corsMiddleware *cors.Cors) gin.HandlerFunc {
	return func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	}
}
