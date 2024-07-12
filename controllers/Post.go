package controllers

import (
	"gelio/m/initializers"
	"gelio/m/middleware"
	"gelio/m/models"

	"github.com/gin-gonic/gin"
)

type Post struct{}

func PostController() *Post {
	return &Post{}
}

func (Post) GetPosts(c *gin.Context) {
	id := c.Param("id")
	offset := c.Param("offset")
	limit := c.Param("limit")

	var Posts []models.Post

	err := initializers.DB.Select(&Posts, `SELECT post_id, body, post.user_id, created_date, image_id
      FROM Message
      INNER JOIN post ON post.user_id = Message.sender_id 
      WHERE receiver_id = $1
      UNION
      SELECT post_id, body, post.user_id, created_date, image_id FROM followers 
      INNER JOIN post on post.user_id = followers.user_id WHERE followers.follower_id = $1
	  UNION SELECT * FROM post where user_id = $1 order by post_id desc limit $2 offset $3`, id, limit, offset)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, Posts)

}

func (Post) UploadPost(c *gin.Context) {
	var body struct {
		Message     string
		UserId      int
		CreatedDate string
		ImageId     int
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := initializers.DB.Exec("insert into post (body, user_id, created_date, image_id) values ($1,$2,$3,$4) returning post_id", body.Message, body.UserId, body.CreatedDate, body.ImageId)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, true)
}

func (p *Post) InitializeRoutes(r *gin.Engine) {
	r.GET("/Posts/:id/:offset/:limit", middleware.RequireAuth, p.GetPosts)
	r.POST("/Post", middleware.RequireAuth, p.UploadPost)
}
