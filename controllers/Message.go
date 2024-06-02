package controllers

import (
	"fmt"
	"gelio/m/initializers"
	"gelio/m/models"

	"github.com/gin-gonic/gin"
)

func LoadMessages(c *gin.Context) {
	id := c.Param("id")

	var Users []models.User

	err := initializers.DB.Select(&Users,
		`SELECT user_id, username, password, created_date, is_active, person_id, profile_image_id 
     FROM Message
     INNER JOIN users ON users.user_id = Message.receiver_id 
     WHERE sender_id = $1
     UNION 
     SELECT user_id, username, password, created_date, is_active, person_id, profile_image_id 
     FROM Message
     INNER JOIN users ON users.user_id = Message.sender_id 
     WHERE receiver_id = $1`, id)

	if err != nil {
		fmt.Println(err)
		c.JSON(400, nil)
		return

	}

	c.JSON(200, Users)
}
