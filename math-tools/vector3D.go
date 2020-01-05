package mathtools

import "strconv"

// Vector3D is a 3 dimensional struct
type Vector3D struct {
	Vector2D
	Z float64
}

// IsZeroVector checks if a vector is a zero vector
func (v Vector3D) IsZeroVector() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0
}

// Length gives the length of the vector
func (v Vector3D) Length() float64 {
	return Distance(v, ZeroVector(3))
}

// Normalize normalizes the vector
func (v Vector3D) Normalize() Vector {
	length := v.Length()
	v.X /= length
	v.Y /= length
	v.Z /= length
	return v
}

// Add adds the vector from current
func (v Vector3D) Add(v2 Vector) Vector {
	return Vector3D{Vector2D{v.X + v2.(Vector3D).X, v.Y + v2.(Vector3D).Y}, v.Z + v2.(Vector3D).Z}
}

// Subtract subtracts the vector from current
func (v Vector3D) Subtract(v2 Vector) Vector {
	return Vector3D{Vector2D{v.X - v2.(Vector3D).X, v.Y - v2.(Vector3D).Y}, v.Z - v2.(Vector3D).Z}
}

// Multiply multplies the vector with a scalar
func (v Vector3D) Multiply(s float64) Vector {
	return Vector3D{Vector2D{v.X * s, v.Y * s}, v.Z * s}
}

// Divide divides the vector with a scalar
func (v Vector3D) Divide(s float64) Vector {
	return Vector3D{Vector2D{v.X / s, v.Y / s}, v.Z / s}
}

// Dot gives the dot product
func (v Vector3D) Dot(v1 Vector) float64 {
	v2 := v1.(Vector3D)
	return v.X * v2.X + v.Y * v2.Y + v.Z * v2.Z
}

// Cross gives the cross product
func (v Vector3D) Cross(v1 Vector) Vector {
	v2 := v1.(Vector3D)
	x := v.Y*v2.Z-v.Z*v2.Y
	y := v.Z*v2.X-v.X*v2.Z
	z := v.X*v2.Y-v.Y*v2.X
	return Vector3D{Vector2D{x, y}, z}
}

// ToString writes the vector out as a string
func (v Vector3D) ToString() string {
	return "(" + strconv.FormatFloat(v.X, 'f', 2, 64) + ", " + strconv.FormatFloat(v.Y, 'f', 2, 64) + ", " + strconv.FormatFloat(v.Z, 'f', 2, 64) + ")"
}

// GetDimension gets the dimension of the vector
func (v Vector3D) GetDimension() int {
	return 3
}
