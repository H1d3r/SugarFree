package Output

import (
	"fmt"
	"os"
)

// Section struct
type Section struct {
	Name    string
	Entropy float64
}

// WriteToFile function
func WriteToFile(sections []Section, filePath string, fileName string, fileSize int64) {
	// Open file to write
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write to the file
	fmt.Fprintf(file, "[+] Analyzing the %s file\n", fileName)
	fmt.Fprintf(file, "[+] File Size: %d bytes\n", fileSize)
	fmt.Fprintln(file, "[+] PE Sections and their entropy:")

	for _, section := range sections {
		// Write the section and its entropy to the file
		fmt.Fprintf(file, "	>>> \"%s\" Scored Entropy Of Value: %.5f\n", section.Name, section.Entropy)
	}
}
