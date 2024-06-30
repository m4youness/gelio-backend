package controllers

import (
	"gelio/m/initializers"
	"gelio/m/models"

	"github.com/gin-gonic/gin"
)

func AddLike(c *gin.Context) {
	var body struct {
		PostId int
		UserId int
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := initializers.DB.Exec("insert into post_likes (post_id, user_id) values ($1, $2)", body.PostId, body.UserId)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, true)

}

func RemoveLike(c *gin.Context) {
	var body struct {
		PostId int
		UserId int
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := initializers.DB.Exec("delete from post_likes where post_id = $1 and user_id = $2", body.PostId, body.UserId)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, true)

}

func IsPostLiked(c *gin.Context) {
	var body struct {
		PostId int
		UserId int
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var PostLike models.PostLikes

	err := initializers.DB.Get(&PostLike, "select * from post_likes where post_id = $1 and user_id = $2", body.PostId, body.UserId)

	if err != nil {
		c.JSON(200, false)
		return
	}

	c.JSON(200, true)

}

func GetAmountOfLikes(c *gin.Context) {
	id := c.Param("id")

	var AmountOfLikes int

	err := initializers.DB.Get(&AmountOfLikes, "select count(*) as post_likes from post_likes group by post_id having post_id = $1", id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(200, 0)
			return
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, AmountOfLikes)

}
