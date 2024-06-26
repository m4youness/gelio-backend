package models

type Image struct {
	ImageId int    `db:"image_id"`
	Url     string `db:"url"`
}
