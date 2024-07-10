package controllers

import (
	"encoding/json"
	"fmt"
	"gelio/m/initializers"
	"gelio/m/middleware"
	"gelio/m/models"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

type Image struct{}

func ImageController() *Image {
	return &Image{}
}

func (Image) UploadImage(c *gin.Context) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File error: " + err.Error()})
		return
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open file: " + err.Error()})
		return
	}
	defer file.Close()

	// Upload the file to Cloudinary
	uploadResult, err := initializers.CloudinaryClient.Upload.Upload(c, file, uploader.UploadParams{
		PublicID: fileHeader.Filename,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary error: " + err.Error()})
		return
	}

	res := initializers.DB.QueryRow("insert into image (url) values ($1) RETURNING image_id", uploadResult.SecureURL)

	var ImageId int

	if err := res.Scan(&ImageId); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, ImageId)

}

func (Image) FindImage(c *gin.Context) {
	id := c.Param("id")

	cachedImage, err := initializers.RedisClient.Get(initializers.Ctx, fmt.Sprintf("image:%s", id)).Result()
	if err == nil {
		var image models.Image
		json.Unmarshal([]byte(cachedImage), &image)
		c.JSON(200, image)
		return
	}

	var image models.Image
	err = initializers.DB.Get(&image, "select * from image where image_id = $1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	imageData, _ := json.Marshal(image)
	err = initializers.RedisClient.Set(initializers.Ctx, fmt.Sprintf("image:%s", id), imageData, time.Hour*24).Err()
	if err != nil {
		fmt.Println("Failed to cache image data:", err)
	}

	c.JSON(http.StatusOK, image)
}

func (i *Image) InitializeRoutes(r *gin.Engine) {
	r.POST("/Image", i.UploadImage)
	r.GET("/Image/:id", middleware.RequireAuth, i.FindImage)
}
