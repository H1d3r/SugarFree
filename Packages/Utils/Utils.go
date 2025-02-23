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
func GetFileSize(filePath string) (float64, error) {
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

	// Convert bytes to KB (divide by 1024)
	sizeInKB := float64(fileInfo.Size()) / 1024.0

	// Return the file size in KB
	return sizeInKB, nil
}

// SplitFileName function
func SplitFileName(filename string) (name, extension string) {
	// Find the last occurrence of "."
	lastDot := strings.LastIndex(filename, ".")

	// If there's no dot or the dot is the first character
	if lastDot <= 0 {
		return filename, ""
	}

	// Split the filename
	name = filename[:lastDot]
	extension = filename[lastDot+1:]

	return name, extension
}

// BuildNewName function
func BuildNewName(name, extension, additionalName string) string {
	// Handle cases where extension might or might not have a dot
	ext := extension
	if extension != "" && !strings.HasPrefix(extension, ".") {
		ext = "." + extension
	}

	// If additionalName is a string containing a number
	value, err := strconv.ParseFloat(additionalName, 64)
	if err != nil {
		// Handle error if the string can't be converted to float
		log.Fatal("Error converting string to float: ", err)
		return ""
	}

	additionalName = fmt.Sprintf("%.5f", value)

	// Build the new name: name_additionalName.extension
	return fmt.Sprintf("%s_%s%s", name, additionalName, ext)
}
