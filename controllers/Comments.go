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

type Comments struct{}

func CommentsController() *Comments {
	return &Comments{}
}

func (Comments) GetComments(c *gin.Context) {
	id := c.Param("id")

	cachedComments, err := initializers.RedisClient.Get(initializers.Ctx, fmt.Sprintf("comments:%s", id)).Result()

	if err == nil {
		var Comments []models.Comments
		json.Unmarshal([]byte(cachedComments), &Comments)
		c.JSON(200, Comments)
		return
	}

	var Comments []models.Comments

	err = initializers.DB.Select(&Comments, "select * from comments where post_id = $1", id)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	commentsData, _ := json.Marshal(Comments)
	err = initializers.RedisClient.Set(initializers.Ctx, fmt.Sprintf("comments:%s", id), commentsData, time.Minute*30).Err()
	if err != nil {
		fmt.Println("Failed to cache comments data:", err)
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

	row := initializers.DB.QueryRow("insert into comments (post_id, user_id, message, created_date) values ($1,$2,$3,$4) returning comment_id", body.PostId, body.UserId, body.Message, body.CreatedDate)

	var CommentId int

	if err := row.Scan(&CommentId); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	cacheKey := fmt.Sprintf("comments:%d", body.PostId)
	cachedComments, err := initializers.RedisClient.Get(initializers.Ctx, cacheKey).Result()
	var comments []models.Comments
	if err == nil {
		json.Unmarshal([]byte(cachedComments), &comments)
	}

	newComment := models.Comments{
		CommentId:   CommentId,
		PostId:      body.PostId,
		UserId:      body.UserId,
		Message:     body.Message,
		CreatedDate: body.CreatedDate,
	}

	comments = append(comments, newComment)

	commentsData, _ := json.Marshal(comments)
	err = initializers.RedisClient.Set(initializers.Ctx, cacheKey, commentsData, time.Minute*30).Err()
	if err != nil {
		fmt.Println("Failed to update cached comments data:", err)
	}

	c.JSON(200, true)

}

func (c *Comments) InitializeRoutes(r *gin.Engine) {
	r.GET("/Comments/:id", middleware.RequireAuth, c.GetComments)
	r.POST("/Comment", middleware.RequireAuth, c.AddComment)
}
