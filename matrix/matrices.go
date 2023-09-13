/*
 * Copyright (C) 2023 by Jason Figge
 */

package matrix

import (
	"math"
)

type Matrix4X4 [4][4]float64

var (
	identity = &Matrix4X4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
)

func Identity() *Matrix4X4 {
	return identity
}

func Translation(x, y, z float64) *Matrix4X4 {
	return &Matrix4X4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{x, y, z, 1},
	}
}

func Projection(aspectRatio, fovRad, near, far float64) *Matrix4X4 {
	return &Matrix4X4{
		{aspectRatio * fovRad, 0, 0, 0},
		{0, fovRad, 0, 0},
		{0, 0, far / (far - near), 1},
		{0, 0, (-far * near) / (far - near), 0},
	}
}

func RotationX(angle float64) *Matrix4X4 {
	return &Matrix4X4{
		{1, 0, 0, 0},
		{0, math.Cos(angle), -math.Sin(angle), 0},
		{0, math.Sin(angle), math.Cos(angle), 0},
		{0, 0, 0, 1},
	}
}
func RotationY(angle float64) *Matrix4X4 {
	return &Matrix4X4{
		{math.Cos(angle), 0, math.Sin(angle), 0},
		{0, 1, 0, 0},
		{-math.Sin(angle), 0, math.Cos(angle), 0},
		{0, 0, 0, 1},
	}
}

func RotationZ(angle float64) *Matrix4X4 {
	return &Matrix4X4{
		{math.Cos(angle), -math.Sin(angle), 0, 0},
		{math.Sin(angle), math.Cos(angle), 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (m *Matrix4X4) Multiply(m1 *Matrix4X4) *Matrix4X4 {
	mo := &Matrix4X4{}
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			mo[r][c] = m[r][0]*m1[0][c] + m[r][1]*m1[1][c] + m[r][2]*m1[2][c] + m[r][3]*m1[3][c]
		}
	}
	return mo
}
