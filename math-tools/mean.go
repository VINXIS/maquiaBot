package mathtools

import "math"

// ArithmeticMean calculates the arithmetic mean for a slice of values
func ArithmeticMean(nums []float64) (mean float64) {
	for _, num := range nums {
		mean += num
	}
	mean /= float64(len(nums))
	return mean
}

// GeometricMean calculates the geometric mean for a slice of values
func GeometricMean(nums []float64) (mean float64) {
	mean = 1.0
	for _, num := range nums {
		mean *= num
	}
	mean = math.Pow(mean, 1/float64(len(nums)))
	return mean
}

// HarmonicMean calculates the harmonic mean for a slice of values
func HarmonicMean(nums []float64) (mean float64) {
	for _, num := range nums {
		mean += 1.0 / num
	}
	mean = float64(len(nums)) / mean
	return mean
}
