package upload

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/davidrbourke/ImageUploader/Backend/utils"
)

// UploadFile api end point handler for uploading an image file
func UploadFile(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)

	fmt.Println("File upload endpoint hit")
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("uploadedFile")
	if err != nil {
		fmt.Println("Error retrieving the file")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded file: %+v\n", handler.Filename)
	fmt.Printf("File size: %+v\n", handler.Size)
	fmt.Printf("MIME header: %+v\n", handler.Header)

	//tempFile, err := ioutil.TempFile("temp-images", handler.Filename)
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("temp-images/"+handler.Filename, buf.Bytes(), 666)

	if err != nil {
		fmt.Println(err)
	}

	ToStorageAccount("temp-images/" + handler.Filename)

	fmt.Fprintf(w, "Successfully uploaded file\n")
}
