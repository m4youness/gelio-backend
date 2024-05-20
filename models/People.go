package models

import "time"

type Person struct {
	person_id     int
	first_name    string
	last_name     string
	gender_id     int
	phone_number  string
	email         string
	date_of_birth time.Time
	country_id    int
}
