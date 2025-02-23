package Reduce

// ApplyStrategy function
func ApplyStrategy(binaryData []byte, numZeroBytes int) []byte {
	zeroBytes := make([]byte, numZeroBytes)

	result := make([]byte, len(binaryData)+numZeroBytes)
	copy(result, binaryData)
	copy(result[len(binaryData):], zeroBytes)

	return result
}
