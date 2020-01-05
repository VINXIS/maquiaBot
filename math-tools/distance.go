package mathtools

import "math"

// Distance calculates the distance between 2 vectors
func Distance(v1, v2 Vector) (distance float64) {
	if v1.GetDimension() == 2 || v2.GetDimension() == 2 {
		v21 := v1.(Vector2D)
		v22 := v2.(Vector2D)
		xDist := math.Pow(v21.X-v22.X, 2.0)
		yDist := math.Pow(v21.Y-v22.Y, 2.0)
		distance = math.Sqrt(xDist + yDist)
	} else {
		v31 := v1.(Vector3D)
		v32 := v2.(Vector3D)
		xDist := math.Pow(v31.X-v32.X, 2.0)
		yDist := math.Pow(v31.Y-v32.Y, 2.0)
		zDist := math.Pow(v31.Z-v32.Z, 2.0)
		distance = math.Sqrt(xDist + yDist + zDist)
	}
	return distance
}
