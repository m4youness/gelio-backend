package controllers

import (
	"gelio/m/initializers"
	"gelio/m/models"

	"github.com/gin-gonic/gin"
)

func GetComments(c *gin.Context) {
	id := c.Param("id")

	var Comments []models.Comments

	err := initializers.DB.Select(&Comments, "select * from comments where post_id = $1", id)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, Comments)

}

func AddComment(c *gin.Context) {
	var body struct {
		PostId      int
		UserId      int
		Message     string
		CreatedDate string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := initializers.DB.Exec("insert into comments (post_id, user_id, message, created_date) values ($1,$2,$3,$4)", body.PostId, body.UserId, body.Message, body.CreatedDate)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, true)
}
