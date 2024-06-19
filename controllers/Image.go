package controllers

import (
<<<<<<< HEAD
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
=======
	"fmt"
	"gelio/m/initializers"
	"github.com/gin-gonic/gin"
)

func AddImage(c *gin.Context) {
	var body struct {
		ImageName string
		ImagePath string
	}

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		return
	}

	var ImageId int

	res := initializers.DB.QueryRow("insert into image (imagename, imagedata) values ($1, $2) returning imageid", body.ImageName, body.ImagePath)

	res.Scan(&ImageId)

	if ImageId == 0 {
		fmt.Println(res.Err())
		return
	}

	c.JSON(200, ImageId)

>>>>>>> c74ed8c47bbcad1fb2db51e22715763bdb190b65
}
