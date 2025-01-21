package WordList

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

// WordCategory represents different types of words for context-aware insertion
type WordCategory int

const (
	CategoryNeutral WordCategory = iota
	CategorySystem
	CategoryData
	CategoryCode
)

// WordEntry represents a word with its category and properties
type WordEntry struct {
	Word     string
	Category WordCategory
	// Entropy impact score (0-1), higher means more entropy reduction
	Impact float64
}

// Initialize word lists by category for context-aware insertion
var (
	systemWords = []WordEntry{
		{"system", CategorySystem, 0.7},
		{"kernel", CategorySystem, 0.8},
		{"driver", CategorySystem, 0.7},
		{"device", CategorySystem, 0.6},
		{"module", CategorySystem, 0.7},
		{"config", CategorySystem, 0.6},
		{"service", CategorySystem, 0.8},
		{"process", CategorySystem, 0.7},
		{"memory", CategorySystem, 0.6},
		{"thread", CategorySystem, 0.7},
	}

	dataWords = []WordEntry{
		{"string", CategoryData, 0.6},
		{"buffer", CategoryData, 0.7},
		{"stream", CategoryData, 0.7},
		{"array", CategoryData, 0.6},
		{"struct", CategoryData, 0.7},
		{"object", CategoryData, 0.6},
		{"value", CategoryData, 0.5},
		{"table", CategoryData, 0.6},
		{"index", CategoryData, 0.6},
		{"field", CategoryData, 0.5},
	}

	codeWords = []WordEntry{
		{"function", CategoryCode, 0.8},
		{"method", CategoryCode, 0.7},
		{"class", CategoryCode, 0.7},
		{"return", CategoryCode, 0.7},
		{"import", CategoryCode, 0.6},
		{"export", CategoryCode, 0.6},
		{"static", CategoryCode, 0.7},
		{"const", CategoryCode, 0.6},
		{"public", CategoryCode, 0.6},
		{"private", CategoryCode, 0.7},
	}

	neutralWords = []WordEntry{
		{"data", CategoryNeutral, 0.5},
		{"info", CategoryNeutral, 0.4},
		{"temp", CategoryNeutral, 0.4},
		{"base", CategoryNeutral, 0.4},
		{"main", CategoryNeutral, 0.4},
		{"file", CategoryNeutral, 0.4},
		{"text", CategoryNeutral, 0.4},
		{"name", CategoryNeutral, 0.4},
		{"type", CategoryNeutral, 0.4},
		{"code", CategoryNeutral, 0.4},
	}

	// Pattern-breaking words specifically designed to reduce entropy
	patternBreakers = []string{
		"0000", "ffff", "aaaa", "5555",
		"3333", "cccc", "9999", "6666",
		"null", "void", "zero", "ones",
		"high", "low", "mid", "end",
	}
)

// Initialize random number generator
func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetContextualWord returns a word appropriate for the context
func GetContextualWord(length int, context []byte) []byte {
	// Analyze context to determine appropriate category
	category := analyzeContext(context)

	// Get word list for category
	var wordList []WordEntry
	switch category {
	case CategorySystem:
		wordList = systemWords
	case CategoryData:
		wordList = dataWords
	case CategoryCode:
		wordList = codeWords
	default:
		wordList = neutralWords
	}

	// Select word with highest impact that fits length
	var selected WordEntry
	maxImpact := -1.0
	for _, entry := range wordList {
		if len(entry.Word) <= length && entry.Impact > maxImpact {
			selected = entry
			maxImpact = entry.Impact
		}
	}

	// If no suitable word found, fall back to neutral
	if maxImpact < 0 {
		return GetWordOfLength(length)
	}

	// Pad word if necessary
	result := make([]byte, length)
	copy(result, selected.Word)
	for i := len(selected.Word); i < length; i++ {
		result[i] = '_'
	}

	return result
}

// GetPatternBreakingWord returns a word designed to break entropy patterns
func GetPatternBreakingWord(length int) []byte {
	if length < 4 {
		return GetWordOfLength(length)
	}

	// Select a pattern breaker
	base := patternBreakers[rand.Intn(len(patternBreakers))]

	// Create result with repeating pattern
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = base[i%len(base)]
	}

	return result
}

// analyzeContext determines the appropriate word category based on surrounding bytes
func analyzeContext(context []byte) WordCategory {
	if len(context) == 0 {
		return CategoryNeutral
	}

	// Convert context to string for pattern matching
	contextStr := string(context)

	// Check for system-related patterns
	if strings.Contains(contextStr, "sys") ||
		strings.Contains(contextStr, "drv") ||
		strings.Contains(contextStr, "dev") {
		return CategorySystem
	}

	// Check for data-related patterns
	if strings.Contains(contextStr, "data") ||
		strings.Contains(contextStr, "buf") ||
		strings.Contains(contextStr, "str") {
		return CategoryData
	}

	// Check for code-related patterns
	if strings.Contains(contextStr, "func") ||
		strings.Contains(contextStr, "class") ||
		strings.Contains(contextStr, "method") {
		return CategoryCode
	}

	return CategoryNeutral
}

// FindWordLocations identifies suitable locations for word insertion while
// ensuring optimal placement for entropy reduction
func FindWordLocations(data []byte, minLength int) []int {
	var locations []int

	// Minimum sequence of zeros required for safe word insertion
	requiredZeros := minLength

	// Scan the data for suitable insertion points
	for i := 0; i < len(data)-requiredZeros; i++ {
		if isValidInsertionPoint(data, i, requiredZeros) {
			locations = append(locations, i)
		}
	}

	// Filter out locations that are too close together
	return optimizeLocations(locations, minLength)
}

// isValidInsertionPoint determines if a location is suitable for word insertion
func isValidInsertionPoint(data []byte, offset int, minZeros int) bool {
	// Check for minimum required zero sequence
	zeroCount := 0
	for i := 0; i < minZeros && offset+i < len(data); i++ {
		if data[offset+i] == 0x00 {
			zeroCount++
		} else {
			break
		}
	}

	// Basic zero sequence check
	if zeroCount < minZeros {
		return false
	}

	// Check surrounding context for unsafe patterns
	if !isSafeContext(data, offset, minZeros) {
		return false
	}

	return true
}

// isSafeContext checks if the surrounding area is safe for word insertion
func isSafeContext(data []byte, offset int, length int) bool {
	// Define safe boundaries for context checking
	start := max(0, offset-length)
	end := min(len(data), offset+length*2)

	// Check for critical patterns that shouldn't be modified
	context := data[start:end]

	// Avoid modifying potential string terminators
	if hasStringTerminator(context) {
		return false
	}

	// Avoid modifying potential pointer or size fields
	if mightBePointer(context) {
		return false
	}

	// Avoid modifying potential array length fields
	if mightBeLength(context) {
		return false
	}

	return true
}

// optimizeLocations filters and optimizes insertion points for best entropy reduction
func optimizeLocations(locations []int, minSpacing int) []int {
	if len(locations) == 0 {
		return locations
	}

	// Sort locations if not already sorted
	if !sort.IntsAreSorted(locations) {
		sort.Ints(locations)
	}

	// Initialize result with first location
	result := []int{locations[0]}
	lastAdded := locations[0]

	// Filter locations that are too close together
	for i := 1; i < len(locations); i++ {
		if locations[i]-lastAdded >= minSpacing*2 {
			result = append(result, locations[i])
			lastAdded = locations[i]
		}
	}

	return result
}

// Helper functions for context analysis

func hasStringTerminator(data []byte) bool {
	// Check for C-style string terminator patterns
	for i := 0; i < len(data)-1; i++ {
		if data[i] == 0x00 && (i == 0 || data[i-1] != 0x00) {
			return true
		}
	}
	return false
}

func mightBePointer(data []byte) bool {
	// Check for common pointer sizes (4 or 8 bytes of zeros)
	if len(data) >= 8 {
		zeroCount := 0
		for i := 0; i < 8; i++ {
			if data[i] == 0x00 {
				zeroCount++
			}
		}
		return zeroCount == 4 || zeroCount == 8
	}
	return false
}

func mightBeLength(data []byte) bool {
	// Check for potential array length fields (4 bytes aligned)
	if len(data) >= 4 {
		allZeros := true
		for i := 0; i < 4; i++ {
			if data[i] != 0x00 {
				allZeros = false
				break
			}
		}
		return allZeros
	}
	return false
}

// GetRandomWord returns a random word from the appropriate category
func GetRandomWord() string {
	// Combine all word categories
	allWords := make([]string, 0)
	for _, entry := range systemWords {
		allWords = append(allWords, entry.Word)
	}
	for _, entry := range dataWords {
		allWords = append(allWords, entry.Word)
	}
	for _, entry := range codeWords {
		allWords = append(allWords, entry.Word)
	}
	for _, entry := range neutralWords {
		allWords = append(allWords, entry.Word)
	}

	return allWords[rand.Intn(len(allWords))]
}

// GetRandomWords returns multiple random words joined by a separator
func GetRandomWords(count int, separator string) string {
	words := make([]string, count)
	for i := 0; i < count; i++ {
		words[i] = GetRandomWord()
	}
	return strings.Join(words, separator)
}

// GetWordOfLength returns a random word padded or truncated to the specified length
func GetWordOfLength(length int) []byte {
	word := GetRandomWord()
	result := make([]byte, length)

	// Copy the word, truncating if necessary
	for i := 0; i < length && i < len(word); i++ {
		result[i] = word[i]
	}

	// Pad with underscores if needed
	for i := len(word); i < length; i++ {
		result[i] = '_'
	}

	return result
}

// GetPaddedWord returns a word padded with a specific byte to reach the desired length
func GetPaddedWord(length int, paddingByte byte) []byte {
	word := GetRandomWord()
	result := make([]byte, length)

	// Copy the word
	copy(result, word)

	// Pad remaining space
	for i := len(word); i < length; i++ {
		result[i] = paddingByte
	}

	return result
}

// Utility functions

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
