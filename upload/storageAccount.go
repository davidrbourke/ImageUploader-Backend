package upload

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/davidrbourke/ImageUploader-Backend/utils"
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
func ToStorageAccount(uploadFilename string, customFileName string) {

	containerURL, ctx := initialiseBlob()
	_, err := containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	handleErrors(err)

	blobURL := containerURL.NewBlockBlobURL(customFileName)
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
func GetAllImageNames(index int, length int, filenameFilter string) (ImagesWrapperResponse, error) {
	containerURL, ctx := initialiseBlob()

	fmt.Println("Listing all blobs in the container")

	qp, accountName := getSAS()

	result := make([]ImageResponse, 0)

	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		handleErrors(err)

		marker = listBlob.NextMarker

		for _, blobInfo := range listBlob.Segment.BlobItems {
			urlToImage := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",
				accountName, "qs-image", blobInfo.Name, qp)

			fmt.Print(" Blob name: " + blobInfo.Name)

			if filenameFilter == "" || strings.HasPrefix(blobInfo.Name, filenameFilter) {
				result = append(result, ImageResponse{ImageName: blobInfo.Name, ImageURL: urlToImage})
			}
		}
	}

	start := (index - 1) * length
	end := start + length
	isLastSet := false
	if end > len(result) {
		end = len(result)
		isLastSet = true
	}
	paged := result[start:end]

	return ImagesWrapperResponse{Images: paged, IsLastSet: isLastSet}, nil
}

func getSAS() (string, string) {
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
		log.Fatal(err)
	}
	sasQueryParams, err := azblob.BlobSASSignatureValues{

		Protocol:      azblob.SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(24 * time.Hour),
		ContainerName: "qs-image",
		BlobName:      "",
		Permissions:   azblob.ContainerSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		log.Fatal(err)
	}

	qp := sasQueryParams.Encode()

	return qp, accountName
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
