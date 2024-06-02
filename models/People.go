package models

type Person struct {
	PersonID    int    `db:"person_id"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	GenderID    int    `db:"gender_id"`
	PhoneNumber string `db:"phone_number"`
	Email       string `db:"email"`
	DateOfBirth string `db:"date_of_birth"`
	CountryID   int    `db:"country_id"`
}
