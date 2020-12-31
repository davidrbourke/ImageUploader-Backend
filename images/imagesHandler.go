package images

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davidrbourke/ImageUploader/Backend/upload"
	"github.com/davidrbourke/ImageUploader/Backend/utils"
)

// GetImages is the API endpoint for handling image requests
func GetImages(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)

	fmt.Println("get images endpoint hit")

	fileNames, err := upload.GetAllImageNames()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(fileNames)
}
