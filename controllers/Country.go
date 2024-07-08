package controllers

import (
	"fmt"
	"gelio/m/initializers"
	"gelio/m/middleware"
	"gelio/m/models"

	"github.com/gin-gonic/gin"
)

type Country struct{}

func CountryController() *Country {
	return &Country{}
}

func (Country) GetAllCountries(c *gin.Context) {

	var Countries []models.Country

	err := initializers.DB.Select(&Countries, "Select * from Country")
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, Countries)

}

func (Country) GetCountryIdWithName(c *gin.Context) {
	var body struct {
		CountryName string
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(400, nil)
		fmt.Println(err)
		return
	}

	var Country models.Country

	Err := initializers.DB.Get(&Country, "select * from Country where country_name = $1", body.CountryName)
	if Err != nil {
		c.JSON(400, nil)
		fmt.Println(Err)
		return
	}

	c.JSON(200, Country.CountryID)
}

func (Country) GetCountryNameWithId(c *gin.Context) {
	id := c.Param("id")

	var Country models.Country
	err := initializers.DB.Get(&Country, "select * from Country where country_id = $1", id)

	if err != nil {
		fmt.Println(err)
		c.JSON(400, nil)
		return
	}

	c.JSON(200, Country.CountryName)
}

func (c *Country) InitializeRoutes(r *gin.Engine) {
	r.GET("/Countries", c.GetAllCountries)
	r.POST("/Get/Country/With/Name", c.GetCountryIdWithName)
	r.GET("/Country/:id", middleware.RequireAuth, c.GetCountryNameWithId)
}
