package controllers

import (
	"gelio/m/initializers"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File error: " + err.Error()})
		return
	}

	// Open the file
	fileHeader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file: " + err.Error()})
		return
	}
	defer fileHeader.Close()

	// Upload the file to Cloudinary
	uploadResult, err := initializers.CloudinaryClient.Upload.Upload(c, file, uploader.UploadParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary error: " + err.Error()})
		return
	}

	// Respond with the secure URL
	c.JSON(http.StatusOK, gin.H{"url": uploadResult.SecureURL})
}
