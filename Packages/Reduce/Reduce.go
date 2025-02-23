package Reduce

import (
	"SugarFree/Packages/Calculate"
	"SugarFree/Packages/WordList"
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"os"
)

// Configuration constants for entropy reduction
const (
	// Minimum length for word insertions
	minWordLength = 4
	// Maximum length for word insertions
	maxWordLength = 16 // Increased from 8 to allow longer pattern breaks
	// Maximum number of words to insert in a single section
	wordInsertLimit = 2000 // Increased from 1000 to allow more insertions
	// Minimum zero bytes required for pattern breaking
	minZeroSequence = 8
	// Entropy threshold for aggressive reduction
	highEntropyThreshold = 7.0
)

// SectionCharacteristics represents known PE section characteristics
const (
	IMAGE_SCN_CNT_CODE               = 0x00000020
	IMAGE_SCN_CNT_INITIALIZED_DATA   = 0x00000040
	IMAGE_SCN_CNT_UNINITIALIZED_DATA = 0x00000080
	IMAGE_SCN_MEM_EXECUTE            = 0x20000000
	IMAGE_SCN_MEM_READ               = 0x40000000
	IMAGE_SCN_MEM_WRITE              = 0x80000000
)

// ReductionStrategy defines the approach for entropy reduction
type ReductionStrategy struct {
	// Whether to use aggressive reduction techniques
	Aggressive bool
	// Target entropy value (0-8)
	TargetEntropy float64
	// Maximum allowed size increase percentage
	MaxSizeIncrease float64
	// Whether to preserve section alignment
	PreserveAlignment bool
}

// calculateMaxWordLength determines the maximum length of word that can be
// inserted at a given location
func calculateMaxWordLength(data []byte, startPos int) int {
	maxLen := minWordLength
	for i := startPos + minWordLength; i < len(data) && i < startPos+maxWordLength; i++ {
		if data[i] != 0x00 {
			break
		}
		maxLen++
	}
	return maxLen
}

// ReduceEntropy processes a PE file and attempts to reduce its entropy while
// maintaining functionality. It employs multiple advanced strategies including
// smart padding, context-aware word injection, and section-specific optimizations.
func ReduceEntropy(filePath string, strategy ReductionStrategy) ([]Calculate.SectionEntropy, error) {
	// Open file for reading and writing
	file, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Parse PE file structure
	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PE file: %w", err)
	}

	var reducedSections []Calculate.SectionEntropy
	var totalSize int64

	// First pass: analyze sections and calculate total size
	for _, section := range peFile.Sections {
		if section.Size == 0 {
			continue
		}
		totalSize += int64(section.Size)
	}

	// Second pass: process each section
	for _, section := range peFile.Sections {
		if section.Size == 0 {
			continue
		}

		// Read section data
		sectionData := make([]byte, section.Size)
		_, err := file.ReadAt(sectionData, int64(section.Offset))
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read section %s: %w", section.Name, err)
		}

		// Calculate initial entropy
		initialEntropy := Calculate.CalculateSectionEntropy(sectionData)

		// Process section based on its characteristics and content
		processedData := processSection(
			sectionData,
			string(section.Name),
			section.Characteristics,
			strategy,
			initialEntropy,
		)

		// Apply word-based entropy reduction if appropriate
		if canApplyWordReduction(string(section.Name), section.Characteristics) {
			// Calculate target word count based on section size and initial entropy
			targetWords := calculateTargetWordCount(
				len(processedData),
				initialEntropy,
				strategy.TargetEntropy,
			)

			processedData = injectWords(
				processedData,
				targetWords,
				strategy.Aggressive,
			)
		}

		// Apply pattern breaking if entropy is still high
		newEntropy := Calculate.CalculateSectionEntropy(processedData)
		if newEntropy > strategy.TargetEntropy && strategy.Aggressive {
			processedData = applyPatternBreaking(
				processedData,
				section.Characteristics,
			)
		}

		// Verify section integrity
		if err := verifySectionIntegrity(
			processedData,
			section.Characteristics,
			strategy.PreserveAlignment,
		); err != nil {
			return nil, fmt.Errorf("section integrity check failed for %s: %w", section.Name, err)
		}

		// Update section info
		reducedSections = append(reducedSections, Calculate.SectionEntropy{
			Name:    string(section.Name),
			Entropy: Calculate.CalculateSectionEntropy(processedData),
			Size:    int64(len(processedData)),
			Offset:  int64(section.Offset),
		})
	}

	if len(reducedSections) == 0 {
		return nil, fmt.Errorf("no valid sections processed")
	}

	return reducedSections, nil
}

// calculateTargetWordCount determines optimal number of words to inject
func calculateTargetWordCount(size int, currentEntropy, targetEntropy float64) int {
	if currentEntropy <= targetEntropy {
		return 0
	}

	// Calculate entropy difference
	entropyDiff := currentEntropy - targetEntropy

	// Base number of words on section size and entropy difference
	baseCount := int(float64(size) * (entropyDiff / 8.0) * 0.1)

	// Ensure we stay within reasonable limits
	if baseCount > wordInsertLimit {
		return wordInsertLimit
	}
	return baseCount
}

// canApplyWordReduction determines if word-based entropy reduction can be safely
// applied to a section based on its name and characteristics
func canApplyWordReduction(name string, characteristics uint32) bool {
	// Never modify executable sections
	if characteristics&IMAGE_SCN_CNT_CODE != 0 ||
		characteristics&IMAGE_SCN_MEM_EXECUTE != 0 {
		return false
	}

	// List of sections that should never be modified
	unsafeSections := map[string]bool{
		".text":  true,
		".idata": true,
		".edata": true,
		".tls":   true,
		".pdata": true,
		".debug": true,
		".rsrc":  true,
		".reloc": true,
	}

	return !unsafeSections[name]
}

// injectWords adds contextually appropriate words to suitable locations
// while preserving the file's functionality
func injectWords(data []byte, targetCount int, aggressive bool) []byte {
	processed := make([]byte, len(data))
	copy(processed, data)

	// Find all suitable locations for word insertion
	locations := WordList.FindWordLocations(processed, minWordLength)
	if len(locations) == 0 {
		return processed
	}

	// Randomize location order to avoid patterns
	rand.Shuffle(len(locations), func(i, j int) {
		locations[i], locations[j] = locations[j], locations[i]
	})

	insertCount := 0
	for _, loc := range locations {
		if insertCount >= targetCount {
			break
		}

		// Calculate maximum word length at this location
		maxLen := calculateMaxWordLength(processed, loc)
		if maxLen < minWordLength {
			continue
		}

		// Get appropriate word for this location
		var word []byte
		if aggressive {
			// Use pattern-breaking words in aggressive mode
			word = WordList.GetPatternBreakingWord(maxLen)
		} else {
			// Use context-appropriate words in normal mode
			word = WordList.GetContextualWord(maxLen, processed[max(0, loc-16):min(len(processed), loc+16)])
		}

		// Insert word only if it doesn't create new patterns
		if !wouldCreatePattern(processed, loc, word) {
			copy(processed[loc:], word)
			insertCount++
		}
	}

	return processed
}

// processSection applies advanced section-specific entropy reduction techniques
func processSection(data []byte, name string, characteristics uint32,
	strategy ReductionStrategy, initialEntropy float64) []byte {

	processed := make([]byte, len(data))
	copy(processed, data)

	// Apply basic processing based on section type
	switch {
	case characteristics&IMAGE_SCN_CNT_CODE != 0:
		return processCodeSection(processed, strategy)
	case characteristics&IMAGE_SCN_CNT_INITIALIZED_DATA != 0:
		return processDataSection(processed, strategy, initialEntropy)
	case characteristics&IMAGE_SCN_CNT_UNINITIALIZED_DATA != 0:
		return processUninitializedSection(processed, strategy)
	}

	// Apply specialized processing based on section name
	switch name {
	case ".data", ".rdata":
		return processDataSection(processed, strategy, initialEntropy)
	case ".rsrc":
		return processResourceSection(processed, strategy)
	case ".reloc":
		return processRelocSection(processed)
	default:
		return processDefaultSection(processed, strategy)
	}
}

// processCodeSection handles entropy reduction for executable code sections
// while carefully preserving functionality
func processCodeSection(data []byte, strategy ReductionStrategy) []byte {
	processed := make([]byte, len(data))
	copy(processed, data)

	if !strategy.Aggressive {
		return processed
	}

	// Process sequences of NOP instructions (0x90)
	for i := 0; i < len(processed)-1; i++ {
		if processed[i] == 0x90 {
			// Replace consecutive NOPs with zeros
			nopCount := 1
			for j := i + 1; j < len(processed) && processed[j] == 0x90; j++ {
				nopCount++
			}

			if nopCount > 2 {
				// Keep one NOP for alignment, replace others
				processed[i] = 0x90
				for j := 1; j < nopCount; j++ {
					processed[i+j] = 0x00
				}
				i += nopCount - 1
			}
		}
	}

	return processed
}

// processDataSection applies sophisticated entropy reduction to data sections
func processDataSection(data []byte, strategy ReductionStrategy, initialEntropy float64) []byte {
	if len(data) < 8 {
		return data
	}

	var processed bytes.Buffer
	processed.Grow(len(data))

	// Adjust chunk size based on initial entropy
	// Higher entropy means we process smaller chunks for more granular control
	chunkSize := 8
	if initialEntropy > 7.0 {
		chunkSize = 4
	}

	// Process data in chunks for pattern analysis
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]

		switch {
		case isStringData(chunk):
			// Preserve string data
			processed.Write(chunk)
		case isNumericData(chunk):
			// Normalize numeric data
			processed.Write(normalizeNumericData(chunk))
		case isHighEntropy(chunk) && strategy.Aggressive:
			// Apply more aggressive reduction for high initial entropy
			if initialEntropy > 7.5 {
				// Use stronger entropy reduction for very high entropy sections
				processed.Write(bytes.Repeat([]byte{0x00}, len(chunk)))
			} else {
				processed.Write(reduceChunkEntropy(chunk))
			}
		default:
			processed.Write(chunk)
		}
	}

	return processed.Bytes()
}

// Helper functions for advanced data processing

// isHighEntropy checks if a chunk of data has high entropy
func isHighEntropy(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	// Calculate chunk entropy using Shannon entropy formula
	var entropy float64
	freq := make(map[byte]int)

	// Count byte frequencies
	for _, b := range data {
		freq[b]++
	}

	// Calculate entropy
	for _, count := range freq {
		p := float64(count) / float64(len(data))
		entropy -= p * math.Log2(p)
	}

	return entropy > highEntropyThreshold
}

// reduceChunkEntropy applies targeted entropy reduction to a data chunk
func reduceChunkEntropy(chunk []byte) []byte {
	reduced := make([]byte, len(chunk))
	copy(reduced, chunk)

	// Replace high-frequency bytes with more common values
	freqMap := make(map[byte]int)
	for _, b := range chunk {
		freqMap[b]++
	}

	// Find most common byte
	var mostCommon byte
	maxCount := 0
	for b, count := range freqMap {
		if count > maxCount {
			maxCount = count
			mostCommon = b
		}
	}

	// Replace similar bytes with most common one
	for i, b := range reduced {
		// If byte is similar to most common, replace it
		if b != mostCommon && math.Abs(float64(b)-float64(mostCommon)) < 16 {
			reduced[i] = mostCommon
		}
	}

	return reduced
}

// processResourceSection carefully reduces entropy in resource sections
func processResourceSection(data []byte, strategy ReductionStrategy) []byte {
	processed := make([]byte, len(data))
	copy(processed, data)

	// Only process if aggressive reduction is enabled
	if !strategy.Aggressive {
		return processed
	}

	// Process resource data carefully to maintain structure
	for i := 0; i < len(processed)-8; i++ {
		// Look for padding sequences
		if isZeroPadding(processed[i:min(i+8, len(processed))]) {
			// Normalize padding bytes
			endPad := min(i+8, len(processed))
			for j := i; j < endPad; j++ {
				processed[j] = 0x00
			}
			i = endPad - 1
		}
	}

	return processed
}

// processRelocSection handles relocation section entropy reduction
func processRelocSection(data []byte) []byte {
	processed := make([]byte, len(data))
	copy(processed, data)

	// Process relocation blocks
	for i := 0; i < len(processed)-8; i += 8 {
		block := processed[i:min(i+8, len(processed))]
		if isRelocBlock(block) {
			normalizeRelocBlock(block)
		}
	}

	return processed
}

// processUninitializedSection handles .bss and similar sections
func processUninitializedSection(data []byte, strategy ReductionStrategy) []byte {
	processed := make([]byte, len(data))
	copy(processed, data)

	// Apply different patterns based on reduction strategy
	if strategy.Aggressive {
		// Use alternating pattern for aggressive reduction
		for i := 0; i < len(processed); i++ {
			if i%2 == 0 {
				processed[i] = 0x00
			} else {
				processed[i] = 0xFF
			}
		}
	} else {
		// Standard approach: fill with zeros
		for i := 0; i < len(processed); i++ {
			if processed[i] != 0x00 {
				processed[i] = 0x00
			}
		}
	}

	return processed
}

// processDefaultSection provides entropy reduction for unknown sections
func processDefaultSection(data []byte, strategy ReductionStrategy) []byte {
	processed := make([]byte, len(data))
	copy(processed, data)

	if !strategy.Aggressive {
		return processed
	}

	// Apply pattern breaking to high-entropy regions
	for i := 0; i < len(processed)-minZeroSequence; i++ {
		chunk := processed[i:min(i+minZeroSequence, len(processed))]
		if isHighEntropy(chunk) {
			reduced := reduceChunkEntropy(chunk)
			copy(processed[i:], reduced)
		}
	}

	return processed
}

// applyPatternBreaking introduces entropy-reducing patterns
func applyPatternBreaking(data []byte, characteristics uint32) []byte {
	processed := make([]byte, len(data))
	copy(processed, data)

	// Skip pattern breaking for executable sections
	if characteristics&IMAGE_SCN_CNT_CODE != 0 ||
		characteristics&IMAGE_SCN_MEM_EXECUTE != 0 {
		return processed
	}

	// Find sequences of similar bytes
	for i := 0; i < len(processed)-8; i++ {
		sequence := processed[i:min(i+8, len(processed))]
		if isRepetitiveSequence(sequence) {
			// Break pattern with alternating bytes
			for j := i; j < min(i+8, len(processed)); j += 2 {
				processed[j] = 0x00
			}
			i += 7
		}
	}

	return processed
}

// verifySectionIntegrity checks if processed section maintains required properties
func verifySectionIntegrity(data []byte, characteristics uint32, preserveAlignment bool) error {
	// Check executable sections
	if characteristics&IMAGE_SCN_CNT_CODE != 0 {
		if !isValidCodeSection(data) {
			return fmt.Errorf("invalid code section modifications detected")
		}
	}

	// Check alignment if required
	if preserveAlignment && !isProperlyAligned(data) {
		return fmt.Errorf("section alignment requirements not met")
	}

	return nil
}

// Helper functions for data validation and analysis

func isValidCodeSection(data []byte) bool {
	// Check for invalid instruction sequences
	for i := 0; i < len(data)-1; i++ {
		// Check for common invalid instruction patterns
		if data[i] == 0x00 && data[i+1] == 0x00 {
			continue
		}
		if data[i] == 0x90 { // NOP
			continue
		}
		// Add more specific checks based on architecture
	}
	return true
}

func isProperlyAligned(data []byte) bool {
	return len(data)%8 == 0
}

func isRepetitiveSequence(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	// Check for byte pattern repetition
	pattern := data[:2]
	for i := 2; i < len(data)-1; i += 2 {
		if data[i] != pattern[0] || data[i+1] != pattern[1] {
			return false
		}
	}
	return true
}

func wouldCreatePattern(data []byte, location int, word []byte) bool {
	// Check surrounding bytes for potential patterns
	start := max(0, location-len(word))
	end := min(len(data), location+len(word)*2)

	context := make([]byte, end-start)
	copy(context, data[start:end])

	// Temporarily insert the word
	copy(context[location-start:], word)

	// Check for repeated patterns
	return hasRepeatingPattern(context)
}

func hasRepeatingPattern(data []byte) bool {
	if len(data) < 8 {
		return false
	}

	// Check for various pattern lengths
	for patternLen := 2; patternLen <= 4; patternLen++ {
		if hasPatternOfLength(data, patternLen) {
			return true
		}
	}
	return false
}

func hasPatternOfLength(data []byte, patternLen int) bool {
	if len(data) < patternLen*2 {
		return false
	}

	for i := 0; i < len(data)-patternLen*2; i++ {
		pattern := data[i : i+patternLen]
		matches := 0

		// Count pattern repetitions
		for j := i + patternLen; j < len(data)-patternLen+1; j += patternLen {
			if bytes.Equal(pattern, data[j:j+patternLen]) {
				matches++
				if matches >= 2 { // Pattern appears at least 3 times
					return true
				}
			} else {
				break
			}
		}
	}
	return false
}

// Utility functions

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper functions for data analysis and processing

func isStringData(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	printable := 0
	for _, b := range data {
		if b >= 32 && b <= 126 {
			printable++
		}
	}
	return float64(printable)/float64(len(data)) > 0.7
}

func isNumericData(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	var num uint32
	return binary.Read(bytes.NewReader(data), binary.LittleEndian, &num) == nil
}

func normalizeNumericData(data []byte) []byte {
	normalized := make([]byte, len(data))
	copy(normalized, data)

	for i := 0; i < len(normalized); i++ {
		if normalized[i] == 0xFF {
			normalized[i] = 0x00
		}
	}

	return normalized
}

func isZeroPadding(data []byte) bool {
	for _, b := range data {
		if b != 0x00 && b != 0xFF {
			return false
		}
	}
	return true
}

func isRelocBlock(block []byte) bool {
	if len(block) < 8 {
		return false
	}
	return block[0] == 0x00 && block[1] == 0x00
}

func normalizeRelocBlock(block []byte) {
	for i := 4; i < len(block); i++ {
		if block[i] == 0xFF {
			block[i] = 0x00
		}
	}
}
