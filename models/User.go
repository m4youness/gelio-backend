package models

import "time"

type User struct {
	UserID       int       `db:"user_id"`
	Username     string    `db:"username"`
	Password     string    `db:"password"`
	CreatedDate  time.Time `db:"created_date"`
	IsActive     bool      `db:"is_active"`
	ProfileImage string    `db:"profile_image"`
	PersonID     int       `db:"person_id"`
}
