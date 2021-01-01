package upload

// ImagesWrapperResponse wraps the images response
type ImagesWrapperResponse struct {
	Images    []ImageResponse
	IsLastSet bool
}
