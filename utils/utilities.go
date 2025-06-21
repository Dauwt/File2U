package utils

import (
	"fmt"
	"os"
)

func GetFileSizeBYTES(fileName string) (int64, error) {
	info, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	fmt.Println("Size in bytes:", info.Size())
	return info.Size(), nil
}

func GetTotalFileSizeBYTES(fileNames []string) (int64, error) {
	var totalSize int64

	for _, fileName := range fileNames {
		fileSize, err := GetFileSizeBYTES(fileName)
		if err != nil {
			return 0, fmt.Errorf("error while getting file size: %d", err)
		}
		totalSize += fileSize
	}

	return totalSize, nil
}
