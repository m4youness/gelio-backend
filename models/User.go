package models

type User struct {
	UserID         int    `db:"user_id"`
	Username       string `db:"username"`
	Password       string `db:"password"`
	CreatedDate    string `db:"created_date"`
	IsActive       bool   `db:"is_active"`
	ProfileImageId *int   `db:"profile_image_id"`
	PersonID       int    `db:"person_id"`
}
