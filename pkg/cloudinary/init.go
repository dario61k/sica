package cloudinary

import (
	"fmt"
	"os"
	"sync"

	"github.com/cloudinary/cloudinary-go/v2"
)

type Cloudinary struct {
	C *cloudinary.Cloudinary
}

var (
	once                sync.Once
	singletonCloudinary *Cloudinary
	initErr             error
)

func GetCloudinary() (*Cloudinary, error) {
	once.Do(func() {
		fmt.Println("CLOUD_NAME:", os.Getenv("C_CLOUD_NAME"))
		fmt.Println("API_KEY:", os.Getenv("C_API_KEY"))
		fmt.Println("API_SECRET:", os.Getenv("C_API_SECRET"))

		cld, err := cloudinary.NewFromParams(
			os.Getenv("C_CLOUD_NAME"),
			os.Getenv("C_API_KEY"),
			os.Getenv("C_API_SECRET"),
		)
		if err != nil {
			initErr = fmt.Errorf("error initializing Cloudinary: %w", err)
			return
		}
		singletonCloudinary = &Cloudinary{C: cld}
	})
	return singletonCloudinary, initErr
}
