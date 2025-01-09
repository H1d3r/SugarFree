package Utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// CheckGoVersio function
func CheckGoVersion() {
	version := runtime.Version()
	version = strings.Replace(version, "go1.", "", -1)
	verNumb, _ := strconv.ParseFloat(version, 64)
	if verNumb < 19.1 {
		logger := log.New(os.Stderr, "[!] ", 0)
		logger.Fatal("The version of Go is to old, please update to version 1.19.1 or later...\n")
	}
}

// GetAbsolutePath function
func GetAbsolutePath(filename string) (string, error) {
	// Get the absolute path of the file
	absolutePath, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}
	return absolutePath, nil
}

// GetFileSize function
func GetFileSize(filePath string) (int64, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Get file information
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %v", err)
	}

	// Return the file size
	return fileInfo.Size(), nil
}
