package main

import (
	"fmt"
	"net/http"

	"github.com/davidrbourke/ImageUploader/Backend/images"
	"github.com/davidrbourke/ImageUploader/Backend/upload"
)

func setupRoutes() {
	http.HandleFunc("/images", images.GetImages)
	http.HandleFunc("/upload", upload.UploadFile)

	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Starting file uploader")
	setupRoutes()
}
