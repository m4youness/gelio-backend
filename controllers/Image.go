package controllers

import (
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

}
