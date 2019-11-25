package tools

import (
	"math"
)

// Levenshtein calculates the difference between 2 messages
func Levenshtein(messageOne, messageTwo string) float64 {
	if len(messageOne) == 0 {
		return 0
	}
	if len(messageTwo) == 0 {
		return 0
	}

	messageOneLength := len(messageOne)
	messageTwoLength := len(messageTwo)

	matrix := make([][]float64, messageOneLength+1)
	for i := range matrix {
		matrix[i] = make([]float64, messageTwoLength+1)
	}

	for i := 1; i <= messageOneLength; i++ {
		matrix[i][0] = float64(i)
	}
	for j := 1; j <= messageTwoLength; j++ {
		matrix[0][j] = float64(j)
	}

	for j := 1; j <= messageTwoLength; j++ {
		for i := 1; i <= messageOneLength; i++ {
			substitutionCost := 0.0
			if messageOne[i-1] != messageTwo[j-1] {
				substitutionCost++
			}
			matrix[i][j] = math.Min(matrix[i-1][j]+1.0, math.Min(matrix[i][j-1]+1.0, matrix[i-1][j-1]+substitutionCost))
		}
	}
	return math.Max(0, matrix[messageOneLength][messageTwoLength])
}
