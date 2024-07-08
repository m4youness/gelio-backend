package initializers

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go"
)

var CloudinaryClient *cloudinary.Cloudinary

func CloudinaryConnect() {
	var err error

	ApiKey := os.Getenv("ClOUDINARY_API_KEY")
	ApiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	CloudinaryClient, err = cloudinary.NewFromParams("geliobackend", ApiKey, ApiSecret)

	if err != nil {
		log.Fatal(err)
		return
	}

}
