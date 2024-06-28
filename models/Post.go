package models

type Post struct {
	PostId      int    `db:"post_id"`
	Body        string `db:"body"`
	UserId      int    `db:"user_id"`
	CreatedDate string `db:"created_date"`
	ImageId     int    `db:"image_id"`
}
