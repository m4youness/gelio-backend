package controllers

import (
	"gelio/m/initializers"
	"gelio/m/models"
	"net/http"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
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

func FindImage(c *gin.Context) {
	id := c.Param("id")

	var Image models.Image

	err := initializers.DB.Get(&Image, "select * from image where image_id = $1", id)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, Image)

}
