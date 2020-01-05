package mathtools

// Vector is the interface for all vector types
type Vector interface {
	IsZeroVector() bool

	Length() float64
	Normalize() Vector

	Add(Vector) Vector
	Subtract(Vector) Vector
	Multiply(s float64) Vector
	Divide(s float64) Vector

	Dot(Vector) float64
	Cross(Vector) Vector

	ToString() string

	GetDimension() int
}

// ZeroVector gives the zero vector
func ZeroVector(dim int) Vector {
	if dim == 2 {
		return Vector2D{}
	}
	return Vector3D{}
}
