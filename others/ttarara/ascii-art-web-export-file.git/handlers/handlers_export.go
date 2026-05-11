package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ExportToFile(data string) (string, error) {

	rootDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current working directory:", err)
		return "", err
	}

	filePath := filepath.Join(rootDir, "temp", "output.txt")

	dir := filepath.Dir(filePath)
	err = os.MkdirAll(dir, 0766) // Create directory if it doesn't exist
	if err != nil {
		log.Println("Error creating directory:", err)
		return "", err
	}

	err = os.WriteFile(filePath, []byte(data), 0766)
	if err != nil {
		log.Println("Error writing to file:", err)
		return "", err
	}

	err = os.Chmod(filePath, 0766)
	if err != nil {
		log.Println("Error setting file permissions:", err)
		return "", err
	}

	fmt.Println("Data successfully written to", filePath)
	return filePath, nil
}
