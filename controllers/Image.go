package controllers

import (
	"gelio/m/initializers"
	"net/http"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File error: " + err.Error()})
		return
	}

	uploadResult, err := initializers.CloudinaryClient.Upload.Upload(c, file, uploader.UploadParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": uploadResult.SecureURL})

}
