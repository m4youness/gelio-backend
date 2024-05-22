package controllers

import (
	"fmt"
	"gelio/m/initializers"

	"github.com/gin-gonic/gin"
)

func AddPerson(c *gin.Context) {
	var body struct {
		FirstName   string
		LastName    string
		GenderID    int
		PhoneNumber string
		Email       string
		DateOfBirth string
		CountryID   int
	}

	err := c.Bind(&body)

	if err != nil {
		fmt.Println(err)
		return
	}

	res := initializers.DB.QueryRow("insert into people (first_name, last_name, gender_id, phone_number, email, date_of_birth, country_id) values ($1, $2, $3, $4, $5, $6, $7) RETURNING person_id", body.FirstName, body.LastName, body.GenderID, body.PhoneNumber, body.Email, body.DateOfBirth, body.CountryID)
	var person_id int
	res.Scan(&person_id)

	if person_id == 0 {
		fmt.Println("Err inserting data in the db")
	}

	c.JSON(200, person_id)
}
