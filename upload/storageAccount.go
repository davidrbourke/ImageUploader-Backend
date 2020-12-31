package upload

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/davidrbourke/ImageUploader/Backend/utils"
	"golang.org/x/net/context"
)

func handleErrors(err error) {
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok { // This error is a Service-specific
			switch serr.ServiceCode() { // Compare serviceCode to ServiceCodeXxx constants
			case azblob.ServiceCodeContainerAlreadyExists:
				fmt.Println("Received 409. Container already exists")
				return
			}
		}
		log.Fatal(err)
	}
}

// ToStorageAccount sends an image file to the azure storage account
func ToStorageAccount(uploadFilename string) {

	containerURL, ctx := initialiseBlob()
	_, err := containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	handleErrors(err)

	blobURL := containerURL.NewBlockBlobURL(uploadFilename)
	file, err := os.Open(uploadFilename)
	handleErrors(err)

	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	handleErrors(err)

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(uploadFilename)
	if err != nil {
		log.Fatal(err)
	}
}

// GetAllImageNames returns a list of all images
func GetAllImageNames() ([]string, error) {
	containerURL, ctx := initialiseBlob()

	fmt.Println("Listing all blobs in the container")

	result := make([]string, 0)

	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		handleErrors(err)

		marker = listBlob.NextMarker

		for _, blobInfo := range listBlob.Segment.BlobItems {
			fmt.Print(" Blob name: " + blobInfo.Name + "\n")
			result = append(result, blobInfo.Name)
		}
	}

	return result, nil
}

func initialiseBlob() (azblob.ContainerURL, context.Context) {
	accountName, err := utils.GetStorageAccountName()
	if err != nil {
		panic(err)
	}

	accountKey, err := utils.GetStorageAccountKey()
	if err != nil {
		panic(err)
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Error creating credential")
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	containerName := "qs-image"

	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	containerURL := azblob.NewContainerURL(*URL, p)
	ctx := context.Background()

	return containerURL, ctx
}
