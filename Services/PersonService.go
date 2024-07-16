package services

import (
	"gelio/m/initializers"
	"gelio/m/models"
)

type PersonService struct{}

func (PersonService) NewPerson(personId int, firstname string, lastname string, genderId int, phonenumber string, email string, dateOfBirth string, countryId int) models.Person {
	return models.Person{
		PersonID:    personId,
		FirstName:   firstname,
		LastName:    lastname,
		GenderID:    genderId,
		PhoneNumber: phonenumber,
		Email:       email,
		DateOfBirth: dateOfBirth,
		CountryID:   countryId,
	}
}

func (PersonService) GetPersonWithId(personId interface{}) (models.Person, error) {
	var Person models.Person
	err := initializers.DB.Get(&Person, "select * from People where person_id = $1", personId)

	if err != nil {
		return models.Person{}, err
	}

	return Person, nil
}

func (PersonService) CreatePerson(person models.Person) (int, error) {
	res := initializers.DB.QueryRow("insert into people (first_name, last_name, gender_id, phone_number, email, date_of_birth, country_id) values ($1, $2, $3, $4, $5, $6, $7) RETURNING person_id", person.FirstName, person.LastName, person.GenderID, person.PhoneNumber, person.Email, person.DateOfBirth, person.CountryID)
	var person_id int
	err := res.Scan(&person_id)

	if err != nil {
		return -1, err
	}

	return person_id, nil

}

func (PersonService) UpdatePerson(person models.Person) error {
	_, err := initializers.DB.Exec("update people set first_name = $1, last_name = $2, email = $3, phone_number = $4, country_id = $5, gender_id = $6 where person_id = $7", person.FirstName, person.LastName, person.Email, person.PhoneNumber, person.CountryID, person.GenderID, person.PersonID)
	if err != nil {
		return err
	}

	return nil
}
