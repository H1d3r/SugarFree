package WordList

import (
	"math/rand"
	"strings"
)

// Predefined list of unique English words (expand as needed)
var englishWords = []string{
	"aaron", "abandoned", "abdomen", "aberdeen", "abilities",
	"ability", "aboriginal", "about", "above", "abraham",
	"absence", "absolute", "absolutely", "absorb", "abstract",
	"absurd", "abuse", "academic", "academy", "accelerate",
	"accent", "acceptable", "acceptance", "access", "accessible",
	"accessories", "accident", "accommodate", "accommodation", "accompany",
	// Add more words as needed
}

// SelectWords function
func SelectWords(numWords int) []string {
	if numWords <= len(englishWords) {
		// Shuffle and take the first numWords
		rand.Shuffle(len(englishWords), func(i, j int) {
			englishWords[i], englishWords[j] = englishWords[j], englishWords[i]
		})
		return englishWords[:numWords]
	}

	// If we need more words than available, generate additional random words
	result := make([]string, len(englishWords))
	copy(result, englishWords)

	// Create a map for faster lookup
	wordMap := make(map[string]bool)
	for _, word := range result {
		wordMap[word] = true
	}

	// Generate additional random words
	chars := "abcdefghijklmnopqrstuvwxyz"
	for len(result) < numWords {
		wordLength := rand.Intn(7) + 4 // Random length between 4 and 10
		var word strings.Builder
		for i := 0; i < wordLength; i++ {
			word.WriteByte(chars[rand.Intn(len(chars))])
		}

		newWord := word.String()
		if !wordMap[newWord] {
			result = append(result, newWord)
			wordMap[newWord] = true
		}
	}

	return result
}
