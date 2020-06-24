package mathtools

import "strconv"

// Vector2D is a 2 dimensional struct
type Vector2D struct {
	X float64
	Y float64
}

// IsZeroVector checks if a vector is a zero vector
func (v Vector2D) IsZeroVector() bool {
	return v.X == 0 && v.Y == 0
}

// Length gives the length of the vector
func (v Vector2D) Length() float64 {
	return Distance(v, ZeroVector(2))
}

// Normalize normalizes the vector
func (v Vector2D) Normalize() Vector {
	length := v.Length()
	v.X /= length
	v.Y /= length
	return v
}

// Add adds the vector from current
func (v Vector2D) Add(v1 Vector) Vector {
	v2 := v1.(Vector2D)
	return Vector2D{v.X + v2.X, v.Y + v2.Y}
}

// Subtract subtracts the vector from current
func (v Vector2D) Subtract(v1 Vector) Vector {
	v2 := v1.(Vector2D)
	return Vector2D{v.X - v2.X, v.Y - v2.Y}
}

// Multiply multplies the vector with a scalar
func (v Vector2D) Multiply(s float64) Vector {
	return Vector2D{v.X * s, v.Y * s}
}

// Divide divides the vector with a scalar
func (v Vector2D) Divide(s float64) Vector {
	return Vector2D{v.X / s, v.Y / s}
}

// Dot gives the dot product
func (v Vector2D) Dot(v1 Vector) float64 {
	v2 := v1.(Vector2D)
	return v.X*v2.X + v.Y*v2.Y
}

// Cross gives the cross product (the value is in the Z coordinate)
func (v Vector2D) Cross(v1 Vector) Vector {
	v2 := v1.(Vector2D)
	return Vector3D{Vector2D{}, v.X*v2.Y - v.Y*v2.X}
}

// ToString writes the vector out as a string
func (v Vector2D) ToString() string {
	return "(" + strconv.FormatFloat(v.X, 'f', 2, 64) + ", " + strconv.FormatFloat(v.Y, 'f', 2, 64) + ")"
}

// GetDimension gets the dimension of the vector
func (v Vector2D) GetDimension() int {
	return 2
}
