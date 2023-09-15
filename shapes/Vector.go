/*
 * Copyright (C) 2023 by Jason Figge
 */

package shapes

import (
	"fmt"
	"math"
)

type Vector struct {
	X float64
	Y float64
	Z float64
	W float64
}

func NewVector(x, y, z float64) *Vector {
	return &Vector{X: x, Y: y, Z: z, W: 1}
}

func NewVectorW(x, y, z, w float64) *Vector {
	return &Vector{X: x, Y: y, Z: z, W: w}
}

type VectorTransformations func(*Vector) *Vector
type XXX func() Vector

func center(x, y float64) VectorTransformations {
	return func(v *Vector) *Vector {
		return &Vector{
			X: (v.X + 1) * x,
			Y: (v.Y + 1) * y,
			Z: v.Z,
		}
	}
}

func project() VectorTransformations {
	return func(v *Vector) *Vector {
		if projectionMatrix == nil {
			return &Vector{X: v.X, Y: v.Y, Z: v.Z}
		}
		projected := v.MatrixMultiply(projectionMatrix)
		if projected.W == 0 {
			return projected
		}
		return projected.Divide(projected.W)
	}
}

func (v *Vector) DotProduct(vArray ...*Vector) float64 {
	switch len(vArray) {
	case 0:
		return v.X*v.X + v.Y*v.Y + v.Z*v.Z
	case 1:
		v1 := vArray[0]
		return v.X*v1.X + v.Y*v1.Y + v.Z*v1.Z
	case 2:
		v1 := vArray[0]
		v2 := vArray[1]
		return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
	}
	panic(fmt.Sprintf("Illegal call to Vector.DotProduct with %d arguments", len(vArray)))
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.DotProduct())
}

func (v *Vector) Normalize() *Vector {
	l := v.Length()
	return &Vector{
		X: v.X / l,
		Y: v.Y / l,
		Z: v.Z / l,
	}
}

func (v *Vector) CrossProduct(v1 *Vector) *Vector {
	return &Vector{
		X: v.Y*v1.Z - v.Z*v1.Y,
		Y: v.Z*v1.X - v.X*v1.Z,
		Z: v.X*v1.Y - v.Y*v1.X,
	}
}

func (v *Vector) Add(v1 *Vector) *Vector {
	return &Vector{
		X: v.X + v1.X,
		Y: v.Y + v1.Y,
		Z: v.Z + v1.Z,
	}
}

func (v *Vector) Subtract(v1 *Vector) *Vector {
	return &Vector{
		X: v.X - v1.X,
		Y: v.Y - v1.Y,
		Z: v.Z - v1.Z,
	}
}

func (v *Vector) Multiply(l float64) *Vector {
	return &Vector{
		X: v.X * l,
		Y: v.Y * l,
		Z: v.Z * l,
	}
}

func (v *Vector) Divide(l float64) *Vector {
	return &Vector{
		X: v.X / l,
		Y: v.Y / l,
		Z: v.Z / l,
	}
}

func (v *Vector) MatrixMultiply(matrix *Matrix4X4) *Vector {
	return &Vector{
		X: v.X*matrix[0][0] + v.Y*matrix[1][0] + v.Z*matrix[2][0] + v.W*matrix[3][0],
		Y: v.X*matrix[0][1] + v.Y*matrix[1][1] + v.Z*matrix[2][1] + v.W*matrix[3][1],
		Z: v.X*matrix[0][2] + v.Y*matrix[1][2] + v.Z*matrix[2][2] + v.W*matrix[3][2],
		W: v.X*matrix[0][3] + v.Y*matrix[1][3] + v.Z*matrix[2][3] + v.W*matrix[3][3],
	}
}

func (v *Vector) Map(v1 *Vector) {
	v.X = v1.X
	v.Y = v1.Y
	v.Z = v1.Z
	v.W = v1.W
}
