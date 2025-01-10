package Manager

import (
	"SugarFree/Packages/Colors"
	"debug/pe"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

// SectionEntropy holds the name and entropy of a PE section
type SectionEntropy struct {
	Name    string
	Entropy float64
}

// CalculateEntropy function
func CalculateEntropy(data []byte) float64 {
	var (
		byteCounts [256]int
		dataLength = len(data)
		entropy    float64
	)

	// Count the occurrences of each byte
	for _, b := range data {
		byteCounts[b]++
	}

	// Calculate the entropy
	for _, count := range byteCounts {
		if count > 0 {
			p := float64(count) / float64(dataLength)
			entropy -= p * math.Log2(p)
		}
	}

	return entropy
}

// ReadSections function
func ReadSections(filePath string) ([]SectionEntropy, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Parse the PE file
	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PE file: %w", err)
	}

	// Create a slice to collect results
	var sectionEntropies []SectionEntropy

	// Iterate through the sections and calculate the entropy of each
	for _, section := range peFile.Sections {
		// Read the raw data of the section
		rawData := make([]byte, section.Size)
		_, err := file.ReadAt(rawData, int64(section.VirtualAddress))
		if err != nil {
			log.Printf("Failed to read section %s: %v\n", section.Name, err)
			continue
		}

		// Calculate entropy of the section
		entropy := CalculateEntropy(rawData)

		// Add the section and its entropy to the result slice
		sectionEntropies = append(sectionEntropies, SectionEntropy{
			Name:    string(section.Name),
			Entropy: entropy,
		})
	}

	return sectionEntropies, nil
}

// ColorNameManager function
func ColorNameManager(section string) string {
	// Declare a variable to store the color of the section name
	var sectionNameColor string

	// Switch statement to color the section names
	switch strings.ToLower(section) {
	case ".text":
		sectionNameColor = Colors.BoldRed(section)
	case ".data":
		sectionNameColor = Colors.BoldGreen(section)
	case ".rdata":
		sectionNameColor = Colors.BoldYellow(section)
	case ".pdata":
		sectionNameColor = Colors.BoldWhite(section)
	case ".bss":
		sectionNameColor = Colors.BoldWhite(section)
	case ".idata":
		sectionNameColor = Colors.BoldMagneta(section)
	case ".edata":
		sectionNameColor = Colors.BoldCyan(section)
	case ".rsrc":
		sectionNameColor = Colors.BoldBlue(section)
	case ".tls":
		sectionNameColor = Colors.BoldRed(section)
	case ".reloc":
		sectionNameColor = Colors.BoldGreen(section)
	case ".debug":
		sectionNameColor = Colors.BoldYellow(section)
	case ".xdata":
		sectionNameColor = Colors.BoldMagneta(section)
	default:
		sectionNameColor = Colors.BoldWhite(section)
	}
	return sectionNameColor
}

// FullEntropy function
func FullEntropy(sections []SectionEntropy) float64 {
	// Initialize a variable to hold the sum of all entropies
	var totalEntropy float64

	// Iterate over each section and sum the entropy values
	for _, section := range sections {
		totalEntropy += section.Entropy
	}

	// Calculate the average entropy if needed (or just return the sum)
	return totalEntropy
}
