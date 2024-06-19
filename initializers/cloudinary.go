package initializers

import (
	"fmt"

	"github.com/cloudinary/cloudinary-go"
)

var CloudinaryClient *cloudinary.Cloudinary

func CloudinaryConnect() {
	var err error

	CloudinaryClient, err = cloudinary.NewFromParams("geliobackend", "816928916411477", "2cCoG-oCk6GphgJRjfgar5ToTkE")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Cloudinary Connected")

}
