package mathtools

import "math"

// StandardDeviation calculates the arithmetic mean for a slice of values
func StandardDeviation(nums []float64, sample bool) (stdDev float64) {
	mean := ArithmeticMean(nums)
	for _, num := range nums {
		stdDev += math.Pow(num-mean, 2.0)
	}
	if sample {
		stdDev /= float64(len(nums) - 1)
	} else {
		stdDev /= float64(len(nums))
	}
	stdDev = math.Sqrt(stdDev)
	return stdDev
}
