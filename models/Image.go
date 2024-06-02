package models

var Image struct {
	ImageId   int    `db:"imageid"`
	ImageName string `db:"imagename"`
	ImagePath string `db:"imagedata"`
}
