package images

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/davidrbourke/ImageUploader-Backend/upload"
	"github.com/davidrbourke/ImageUploader-Backend/utils"
)

// GetImages is the API endpoint for handling image requests
func GetImages(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)

	fmt.Println("get images endpoint hit")

	index := r.URL.Query().Get("index")
	length := r.URL.Query().Get("length")
	filenameFilter := r.URL.Query().Get("filter")

	iIndex, err := strconv.Atoi(index)
	if err != nil {
		iIndex = 1
	}

	iLength, err := strconv.Atoi(length)
	if err != nil {
		iLength = 10
	}

	fmt.Printf("%s, %s", index, length)

	imageResponse, err := upload.GetAllImageNames(iIndex, iLength, filenameFilter)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(imageResponse)
}
