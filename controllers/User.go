package controllers

import (
	"fmt"
	"gelio/m/initializers"
	"gelio/m/models"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var body struct {
		Username string
	}

	c.Bind(&body)

	var User models.User

	err := initializers.DB.Get(&User, "select * from users where username = $1", body.Username)
	fmt.Println(body.Username)

	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, User)

}
