package initializers

import (
	"fmt"
	"github.com/cloudinary/cloudinary-go"
	"os"
)

var CloudinaryClient *cloudinary.Cloudinary

func CloudinaryConnect() {
	var err error

	ApiKey := os.Getenv("ClOUDINARY_API_KEY")
	ApiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	CloudinaryClient, err = cloudinary.NewFromParams("geliobackend", ApiKey, ApiSecret)

	if err != nil {
		fmt.Println(err)
		return
	}

}
