package Colors

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	// Bold Colors
	BoldBlue    = color.New(color.FgBlue, color.Bold).SprintFunc()
	BoldRed     = color.New(color.FgRed, color.Bold).SprintFunc()
	BoldGreen   = color.New(color.FgGreen, color.Bold).SprintFunc()
	BoldYellow  = color.New(color.FgYellow, color.Bold).SprintFunc()
	BoldWhite   = color.New(color.FgHiWhite, color.Bold).SprintFunc()
	BoldMagenta = color.New(color.FgMagenta, color.Bold).SprintFunc()
	BoldCyan    = color.New(color.FgCyan, color.Bold).SprintFunc()
)

// Define a slice containing all available color functions
var allColors = []func(a ...interface{}) string{
	BoldBlue, BoldRed, BoldGreen, BoldYellow, BoldWhite, BoldMagenta, BoldCyan,
}

// RandomColor function
// RandomColor selects a random color function from the available ones
func RandomColor() func(a ...interface{}) string {
	rand.Seed(time.Now().UnixNano())
	return allColors[rand.Intn(len(allColors))]
}

// ColorNameManager function
func ColorNameManager(section string) string {
	// Determine color based on section name
	switch strings.ToLower(section) {
	case ".text":
		return BoldYellow(section) // Code section
	case ".data":
		return BoldMagenta(section) // Data section
	case ".rdata":
		return BoldBlue(section) // Read-only data
	case ".pdata":
		return BoldCyan(section) // Exception handling data
	case ".bss":
		return BoldWhite(section) // Uninitialized data
	case ".idata":
		return BoldYellow(section) // Import directory
	case ".edata":
		return BoldMagenta(section) // Export directory
	case ".rsrc":
		return BoldBlue(section) // Resources
	case ".tls":
		return BoldCyan(section) // Thread-local storage
	case ".reloc":
		return BoldWhite(section) // Relocations
	case ".debug":
		return BoldYellow(section) // Debug information
	case ".xdata":
		return BoldMagenta(section) // Exception information
	default:
		return BoldBlue(section) // Unknown sections
	}
}

// CalculateColor2Entropy function
func CalculateColor2Entropy(entropy float64) string {
	// Check if the entropy is less than 5.0
	if entropy < 5.0 {
		return BoldGreen(fmt.Sprintf("%.5f", entropy))
	}

	return BoldRed(fmt.Sprintf("%.5f", entropy))
}
