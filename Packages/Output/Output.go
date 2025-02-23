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

// Write2File function
func Write2File(sections []Section, filePath string, fileName string, fileSize float64, fullEntropy float64, getDateTime string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	WriteBasicInfo(file, fileName, fileSize, fullEntropy, getDateTime)
	WriteSectionInfo(file, sections)
}

// writeBasicInfo writes basic file information
func WriteBasicInfo(file *os.File, fileName string, fileSize float64, entropy float64, getDateTime string) {
	fmt.Fprintf(file, "PE Analysis Report - %s\n\n", getDateTime)
	fmt.Fprintf(file, "File Name: %s\n", fileName)
	fmt.Fprintf(file, "File Size: %f bytes\n", fileSize)
	fmt.Fprintf(file, "Overall PE Entropy: %.5f\n", entropy)
}

// writeSectionInfo writes detailed section information
func WriteSectionInfo(file *os.File, sections []Section) {
	fmt.Fprintln(file, "\nPE Sections Entropy:")
	for _, section := range sections {
		fmt.Fprintf(file, "  >>> \"%s\" Entropy: %.5f\n", section.Name, section.Entropy)
	}
}
