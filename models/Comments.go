package models

type Comments struct {
	CommentId   int    `db:"comment_id"`
	PostId      int    `db:"post_id"`
	UserId      int    `db:"user_id"`
	Message     string `db:"message"`
	CreatedDate string `db:"created_date"`
}
