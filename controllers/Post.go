package controllers

import (
	"gelio/m/initializers"
	"gelio/m/models"

	"github.com/gin-gonic/gin"
)

func GetPosts(c *gin.Context) {
	id := c.Param("id")

	var Posts []models.Post

	err := initializers.DB.Select(&Posts, `SELECT post_id, body, post.user_id, created_date, image_id
      FROM Message
      INNER JOIN post ON post.user_id = Message.sender_id 
      WHERE receiver_id = $1
      UNION
      SELECT post_id, body, post.user_id, created_date, image_id FROM followers 
      INNER JOIN post on post.user_id = followers.user_id WHERE followers.follower_id = $1
	  UNION SELECT * FROM post where user_id = $1`, id)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, Posts)

}

func UploadPost(c *gin.Context) {
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

	res := initializers.DB.QueryRow("insert into post (body, user_id, created_date, image_id) values ($1,$2,$3,$4) returning post_id", body.Message, body.UserId, body.CreatedDate, body.ImageId)

	var PostId int

	err := res.Scan(&PostId)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, PostId)

}
