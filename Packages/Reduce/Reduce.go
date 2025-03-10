package Reduce

import (
	"SugarFree/Packages/WordList"
	"log"
	"strings"
)

// ApplyStrategy function
func ApplyStrategy(binaryData []byte, number int, strategy string) []byte {
	switch strings.ToLower(strategy) {
	case "zero":
		zeroBytes := make([]byte, number)
		//fmt.Println(zeroBytes)

		result := make([]byte, len(binaryData)+number)
		copy(result, binaryData)
		copy(result[len(binaryData):], zeroBytes)

		return result
	case "word":
		// Call function named SelectWords
		words := WordList.SelectWords(number)
		//fmt.Print(words)

		wordsBytes := []byte(strings.Join(words, ""))
		//fmt.Print(wordsBytes)

		result := append(binaryData, wordsBytes...)

		return result
	default:
		log.Fatal("Error: Invalid strategy provided. Please provide a valid strategy to continue...\n\n")
		return nil
	}
}
