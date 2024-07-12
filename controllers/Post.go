package controllers

import (
	"encoding/json"
	"fmt"
	"gelio/m/initializers"
	"gelio/m/middleware"
	"gelio/m/models"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct{}

func PostController() *Post {
	return &Post{}
}

func (Post) GetPosts(c *gin.Context) {
	id := c.Param("id")

	cachedPosts, err := initializers.RedisClient.Get(initializers.Ctx, fmt.Sprintf("posts:%s", id)).Result()

	if err == nil {
		var Posts []models.Post
		json.Unmarshal([]byte(cachedPosts), &Posts)
		c.JSON(200, Posts)
		return
	}

	var Posts []models.Post

	err = initializers.DB.Select(&Posts, `SELECT post_id, body, post.user_id, created_date, image_id
      FROM Message
      INNER JOIN post ON post.user_id = Message.sender_id 
      WHERE receiver_id = $1
      UNION
      SELECT post_id, body, post.user_id, created_date, image_id FROM followers 
      INNER JOIN post on post.user_id = followers.user_id WHERE followers.follower_id = $1
	  UNION SELECT * FROM post where user_id = $1 order by post_id desc`, id)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	postsData, _ := json.Marshal(Posts)
	err = initializers.RedisClient.Set(initializers.Ctx, fmt.Sprintf("posts:%s", id), postsData, time.Minute*30).Err()
	if err != nil {
		fmt.Println("Failed to cache posts data:", err)
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

	res := initializers.DB.QueryRow("insert into post (body, user_id, created_date, image_id) values ($1,$2,$3,$4) returning post_id", body.Message, body.UserId, body.CreatedDate, body.ImageId)

	var PostId int

	err := res.Scan(&PostId)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	cacheKey := fmt.Sprintf("posts:%d", body.UserId)
	cachedPosts, err := initializers.RedisClient.Get(initializers.Ctx, cacheKey).Result()
	var posts []models.Post
	if err == nil {
		json.Unmarshal([]byte(cachedPosts), &posts)
	}

	newPost := models.Post{
		PostId:      PostId,
		Body:        body.Message,
		UserId:      body.UserId,
		CreatedDate: body.CreatedDate,
		ImageId:     body.ImageId,
	}

	// Append the new post to the slice
	posts = append(posts, newPost)

	// Marshal the updated posts slice back into JSON
	postsData, _ := json.Marshal(posts)
	err = initializers.RedisClient.Set(initializers.Ctx, cacheKey, postsData, time.Minute*30).Err()
	if err != nil {
		fmt.Println("Failed to update cached posts data:", err)
	}

	c.JSON(200, PostId)
}

func (p *Post) InitializeRoutes(r *gin.Engine) {
	r.GET("/Posts/:id", middleware.RequireAuth, p.GetPosts)
	r.POST("/Post", middleware.RequireAuth, p.UploadPost)
}
