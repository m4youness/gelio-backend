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
	r.GET("/IsAuthenticated", middleware.RequireAuth, controllers.IsLoggedIn)
	r.GET("/Logout", middleware.RequireAuth, controllers.Logout)
<<<<<<< HEAD
	r.GET("/User/:id", middleware.RequireAuth, controllers.GetUser)
	r.POST("/User/Exists", controllers.DoesUserExist)
	r.GET("/UserId", middleware.RequireAuth, controllers.GetUserId)

	// People
	r.POST("/Person", controllers.AddPerson)
	r.GET("/Person/:id", middleware.RequireAuth, controllers.GetPerson)
=======
	r.GET("/User/:id", controllers.GetUser)
	r.POST("/User/Exists", controllers.DoesUserExist)
	r.GET("/UserId", controllers.GetUserId)

	// People
	r.POST("/Person", controllers.AddPerson)
	r.GET("/Person/:id", controllers.GetPerson)
>>>>>>> c74ed8c47bbcad1fb2db51e22715763bdb190b65

	// Country
	r.GET("/Countries", controllers.GetAllCountries)
	r.POST("/GetCountryWithName", controllers.GetCountryIdWithName)
<<<<<<< HEAD
	r.GET("/Country/:id", middleware.RequireAuth, controllers.GetCountryNameWithId)

	// Image
	r.POST("/Image", middleware.RequireAuth, controllers.UploadImage)

	// Message
	r.GET("/LoadContacts/:id", middleware.RequireAuth, controllers.LoadContacts)
	r.POST("/LoadMessages", middleware.RequireAuth, controllers.LoadMessages)
	r.POST("/Message", middleware.RequireAuth, controllers.SendMessage)
	r.POST("/Contact", middleware.RequireAuth, controllers.AddContact)
	r.POST("/ContactExist", middleware.RequireAuth, controllers.IsPersonNotContact)
=======
	r.GET("/Country/:id", controllers.GetCountryNameWithId)

	// Image
	r.POST("/Image", controllers.AddImage)

	// Message
	r.GET("/LoadContacts/:id", controllers.LoadContacts)
	r.POST("/LoadMessages", controllers.LoadMessages)
	r.GET("/MessageInfo/:id", controllers.GetMessageInfoFromId)
	r.POST("/Message", controllers.SendMessage)
	r.POST("/Contact", controllers.AddContact)
>>>>>>> c74ed8c47bbcad1fb2db51e22715763bdb190b65

	r.Run()

}

func corsHandler(corsMiddleware *cors.Cors) gin.HandlerFunc {
	return func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	}
}
