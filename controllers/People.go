package controllers

import (
	"fmt"
	"gelio/m/initializers"
	"gelio/m/middleware"
	"gelio/m/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type Person struct{}

func PeopleController() *Person {
	return &Person{}
}

func (Person) AddPerson(c *gin.Context) {
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
		fmt.Println(res.Err())
		fmt.Println("Err inserting data in the db")
	}

	c.JSON(200, person_id)
}

func (Person) GetPerson(c *gin.Context) {
	id := c.Param("id")

	var Person models.Person

	err := initializers.DB.Get(&Person, "select * from People where person_id = $1", id)

	dateParts := strings.Split(Person.DateOfBirth, "T")
	Person.DateOfBirth = dateParts[0]

	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, Person)

}

func (p *Person) InitializeRoutes(r *gin.Engine) {
	r.POST("/Person", p.AddPerson)
	r.GET("/Person/:id", middleware.RequireAuth, p.GetPerson)
}
