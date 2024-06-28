package models

type PostLikes struct {
	PostLikesId int `db:"post_likes_id"`
	PostId      int `db:"post_id"`
	UserId      int `db:"user_id"`
}
