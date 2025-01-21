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

// SectionEntropy holds comprehensive information about a PE section's entropy
type SectionEntropy struct {
	Name    string  // Section name
	Entropy float64 // Calculated entropy value
	Size    int64   // Size of the section in bytes
	Offset  int64   // File offset of the section
}

// CalculateEntropy computes Shannon entropy for a byte sequence
// Returns a value between 0 and 8 bits (0 to 8 bits per byte)
func CalculateEntropy(data []byte) float64 {
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

// ReadSections reads and analyzes all sections of a PE file
// Returns entropy information for each valid section
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

		// Calculate entropy for the section
		entropy := CalculateEntropy(rawData[:n])

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

// FullEntropy calculates the weighted average entropy across all sections
// This provides a more accurate representation of the file's overall entropy
func FullEntropy(sections []SectionEntropy) float64 {
	var totalSize int64
	var weightedEntropy float64

	// Calculate weighted sum of entropies
	for _, section := range sections {
		totalSize += section.Size
		weightedEntropy += float64(section.Size) * section.Entropy
	}

	// Handle empty file case
	if totalSize == 0 {
		log.Printf("Warning: No valid sections found for entropy calculation")
		return 0
	}

	// Calculate weighted average
	averageEntropy := weightedEntropy / float64(totalSize)

	// Validate final result
	if averageEntropy < 0 || averageEntropy > 8 {
		log.Printf("Warning: Invalid weighted average entropy: %f, adjusting to valid range", averageEntropy)
		return math.Min(8, math.Max(0, averageEntropy))
	}

	return averageEntropy
}

// ColorNameManager assigns consistent colors to different PE section names
// This helps in visual identification of different section types
func ColorNameManager(section string) string {
	// Determine color based on section name
	switch strings.ToLower(section) {
	case ".text":
		return Colors.BoldRed(section) // Code section
	case ".data":
		return Colors.BoldGreen(section) // Data section
	case ".rdata":
		return Colors.BoldYellow(section) // Read-only data
	case ".pdata":
		return Colors.BoldWhite(section) // Exception handling data
	case ".bss":
		return Colors.BoldWhite(section) // Uninitialized data
	case ".idata":
		return Colors.BoldMagneta(section) // Import directory
	case ".edata":
		return Colors.BoldCyan(section) // Export directory
	case ".rsrc":
		return Colors.BoldBlue(section) // Resources
	case ".tls":
		return Colors.BoldRed(section) // Thread-local storage
	case ".reloc":
		return Colors.BoldGreen(section) // Relocations
	case ".debug":
		return Colors.BoldYellow(section) // Debug information
	case ".xdata":
		return Colors.BoldMagneta(section) // Exception information
	default:
		return Colors.BoldWhite(section) // Unknown sections
	}
}

// ValidateSection performs basic validation of section data
func ValidateSection(section SectionEntropy) error {
	if section.Size < 0 {
		return fmt.Errorf("invalid negative section size: %d", section.Size)
	}
	if section.Offset < 0 {
		return fmt.Errorf("invalid negative section offset: %d", section.Offset)
	}
	if section.Entropy < 0 || section.Entropy > 8 {
		return fmt.Errorf("entropy value out of valid range [0,8]: %f", section.Entropy)
	}
	return nil
}
