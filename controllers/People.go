package controllers

import (
	"fmt"
	"gelio/m/IServices"
	"gelio/m/middleware"
	"gelio/m/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type PersonController struct {
	_IPersonService IServices.IPersonService
}

func NewPersonController(IPersonService IServices.IPersonService) *PersonController {
	return &PersonController{
		_IPersonService: IPersonService,
	}
}

func (u *PersonController) AddPerson(c *gin.Context) {
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

	person := u._IPersonService.NewPerson(-1, body.FirstName, body.LastName, body.GenderID, body.PhoneNumber, body.Email, body.DateOfBirth, body.CountryID)

	personId, err := u._IPersonService.CreatePerson(person)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, personId)
}

func (p *PersonController) GetPerson(c *gin.Context) {
	id := c.Param("id")

	var Person models.Person

	Person, err := p._IPersonService.GetPersonWithId(id)

	dateParts := strings.Split(Person.DateOfBirth, "T")
	Person.DateOfBirth = dateParts[0]

	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, Person)

}

func (p *PersonController) InitializeRoutes(r *gin.Engine) {
	r.POST("/Person", p.AddPerson)
	r.GET("/Person/:id", middleware.RequireAuth, p.GetPerson)
}
