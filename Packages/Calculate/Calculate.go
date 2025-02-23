package Calculate

import (
	"debug/pe"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

// SectionEntropy struct
type SectionEntropy struct {
	Name    string  // Section name
	Entropy float64 // Calculated entropy value
	Size    int64   // Size of the section in bytes
	Offset  int64   // File offset of the section
}

// CalculateFullEntropy function
func CalculateFullEntropy(buffer []byte) float64 {
	entropy := 0.0

	// Create a map to count byte occurrences
	counts := make(map[byte]int)
	for _, b := range buffer {
		counts[b]++
	}

	// Calculate entropy using Shannon's formula
	bufferLen := float64(len(buffer))
	for _, count := range counts {
		p := float64(count) / bufferLen
		entropy += -p * math.Log2(p)
	}

	return entropy
}

// CalculateSectionEntropy function
func CalculateSectionEntropy(data []byte) float64 {
	// Handle empty data case
	if len(data) == 0 {
		return 0
	}

	var (
		byteCounts [256]int
		dataLength = float64(len(data))
		entropy    float64
	)

	// Count occurrences of each byte value
	for _, b := range data {
		byteCounts[b]++
	}

	// Calculate Shannon entropy
	for _, count := range byteCounts {
		if count > 0 {
			probability := float64(count) / dataLength
			entropy -= probability * math.Log2(probability)
		}
	}

	// Ensure entropy stays within valid range
	if entropy < 0 {
		log.Printf("Warning: Negative entropy calculated, adjusting to 0")
		return 0
	}
	if entropy > 8 {
		log.Printf("Warning: Entropy exceeded maximum possible value, adjusting to 8")
		return 8
	}

	return entropy
}

// ReadSections function
func ReadSections(filePath string) ([]SectionEntropy, error) {
	// Open the PE file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Parse the PE file structure
	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PE file: %w", err)
	}

	var sectionEntropies []SectionEntropy

	// Process each section
	for _, section := range peFile.Sections {
		// Skip sections with zero size
		if section.Size == 0 {
			log.Printf("Warning: Skipping zero-size section %s\n", section.Name)
			continue
		}

		// Allocate buffer for section data
		rawData := make([]byte, section.Size)

		// Read section data using correct file offset
		n, err := file.ReadAt(rawData, int64(section.Offset))
		if err != nil {
			// Special handling for .reloc section which often has special requirements
			if strings.EqualFold(section.Name, ".reloc") {
				log.Printf("Warning: Incomplete read of .reloc section: %v\n", err)
				if n == 0 {
					continue
				}
				// Use partial data if some was read
				rawData = rawData[:n]
			} else {
				// For other sections, report the error
				return nil, fmt.Errorf("failed to read section %s: %w", section.Name, err)
			}
		}

		// Call function named CalculateSectionEntropy
		entropy := CalculateSectionEntropy(rawData[:n])

		// Store section information
		sectionEntropies = append(sectionEntropies, SectionEntropy{
			Name:    string(section.Name),
			Entropy: entropy,
			Size:    int64(n),
			Offset:  int64(section.Offset),
		})
	}

	return sectionEntropies, nil
}
