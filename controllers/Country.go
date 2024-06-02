package controllers

import (
	"fmt"
	"gelio/m/initializers"
	"gelio/m/models"

	"github.com/gin-gonic/gin"
)

func GetAllCountries(c *gin.Context) {

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

func GetCountryIdWithName(c *gin.Context) {
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

func GetCountryNameWithId(c *gin.Context) {
	id := c.Param("id")

	var Country models.Country
	err := initializers.DB.Get(&Country, "select * from Country where country_id = $1", id)

	if err != nil {
		fmt.Println(err)
		c.JSON(400, nil)
		return
	}
	fmt.Println(Country.CountryName)
	c.JSON(200, Country.CountryName)
}
