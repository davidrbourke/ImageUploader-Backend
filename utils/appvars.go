package utils

import (
	"fmt"
	"os"
)

// GetStorageAccountName returns the image storage account name from the environment
func GetStorageAccountName() (string, error) {
	saName := os.Getenv("IMAGE_STORAGEACCOUNT_NAME")
	if saName == "" {
		return "", fmt.Errorf("IMAGE_STORAGEACCOUNT_NAME is empty")
	}

	return saName, nil
}

// GetStorageAccountKey returns the image storage account key from the environment
func GetStorageAccountKey() (string, error) {
	saKey := os.Getenv("IMAGE_STORAGEACCOUNT_KEY")
	if saKey == "" {
		return "", fmt.Errorf("IMAGE_STORAGEACCOUNT_KEY is empty")
	}

	return saKey, nil
}
