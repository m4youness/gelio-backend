package models

type Follower struct {
	UserId     int `db:"user_id"`
	FollowerId int `db:"follower_id"`
}
