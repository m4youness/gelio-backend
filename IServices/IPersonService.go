package IServices

import "gelio/m/models"

type IPersonService interface {
	GetPersonWithId(personId interface{}) (models.Person, error)
	CreatePerson(person models.Person) (int, error)
	UpdatePerson(person models.Person) error
	NewPerson(perosonId int, firstname string, lastname string, genderId int, phonenumber string, email string, dateOfBirth string, countryId int) models.Person
}
