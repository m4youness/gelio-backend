package controllers

import (
	"gelio/m/initializers"
	"gelio/m/middleware"
	"gelio/m/models"

	"github.com/gin-gonic/gin"
)

type Comments struct{}

func CommentsController() *Comments {
	return &Comments{}
}

func (Comments) GetComments(c *gin.Context) {
	id := c.Param("id")

	var Comments []models.Comments

	err := initializers.DB.Select(&Comments, "select * from comments where post_id = $1", id)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, Comments)

}

func (Comments) AddComment(c *gin.Context) {
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

func (c *Comments) InitializeRoutes(r *gin.Engine) {
	r.GET("/Comments/:id", middleware.RequireAuth, c.GetComments)
	r.POST("/Comment", middleware.RequireAuth, c.AddComment)
}
