package Output

import (
	"fmt"
	"os"
	"time"
)

// Section struct defines the structure for section information
type Section struct {
	Name    string
	Entropy float64
}

// WriteToFile writes basic section information to a file
func WriteToFile(sections []Section, filePath string, fileName string, fileSize int64, fullEntropy float64) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	writeBasicInfo(file, fileName, fileSize, fullEntropy)
	writeSectionInfo(file, sections)
}

// WriteEntropyReport writes a comprehensive entropy analysis report
func WriteEntropyReport(sections []Section, filePath string, fileName string, fileSize int64,
	initialEntropy float64, finalEntropy float64) {

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write report header
	fmt.Fprintf(file, "=== SugarFree Entropy Analysis Report ===\n")
	fmt.Fprintf(file, "Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Write file information
	writeBasicInfo(file, fileName, fileSize, initialEntropy)

	// Write entropy reduction results
	fmt.Fprintf(file, "\n=== Entropy Reduction Results ===\n")
	fmt.Fprintf(file, "Initial Entropy: %.5f\n", initialEntropy)
	fmt.Fprintf(file, "Final Entropy: %.5f\n", finalEntropy)

	// Calculate and write reduction percentage
	reductionPercent := ((initialEntropy - finalEntropy) / initialEntropy) * 100
	fmt.Fprintf(file, "Entropy Reduction: %.2f%%\n", reductionPercent)

	// Write detailed section information
	fmt.Fprintf(file, "\n=== Section Details ===\n")
	writeSectionInfo(file, sections)

	// Write analysis summary
	writeAnalysisSummary(file, sections)
}

// writeBasicInfo writes basic file information
func writeBasicInfo(file *os.File, fileName string, fileSize int64, entropy float64) {
	fmt.Fprintf(file, "File Name: %s\n", fileName)
	fmt.Fprintf(file, "File Size: %d bytes\n", fileSize)
	fmt.Fprintf(file, "Full PE Entropy: %.5f\n", entropy)
}

// writeSectionInfo writes detailed section information
func writeSectionInfo(file *os.File, sections []Section) {
	fmt.Fprintln(file, "\nPE Sections and their entropy:")
	for _, section := range sections {
		fmt.Fprintf(file, "  >>> \"%s\" Entropy: %.5f\n", section.Name, section.Entropy)
	}
}

// writeAnalysisSummary writes an analysis summary with recommendations
func writeAnalysisSummary(file *os.File, sections []Section) {
	fmt.Fprintf(file, "\n=== Analysis Summary ===\n")

	// Analyze high entropy sections
	highEntropySections := 0
	for _, section := range sections {
		if section.Entropy > 6.5 {
			highEntropySections++
		}
	}

	if highEntropySections > 0 {
		fmt.Fprintf(file, "Found %d section(s) with high entropy (>6.5)\n", highEntropySections)
		fmt.Fprintf(file, "Recommendation: Consider additional entropy reduction techniques for these sections\n")
	} else {
		fmt.Fprintf(file, "All sections have acceptable entropy levels\n")
	}
}
