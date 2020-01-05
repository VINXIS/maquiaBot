package mathtools

// Direction calculates the direction of a vector or the distance between 2 vectors (v1 - v2)
func Direction(v1, v2 Vector) Vector {
	if v2.IsZeroVector() {
		return v1.Normalize()
	}

	v3 := v1.Subtract(v2)
	return v3.Normalize()
}
