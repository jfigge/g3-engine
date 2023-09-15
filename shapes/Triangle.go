/*
 * Copyright (C) 2023 by Jason Figge
 */

package shapes

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Triangle struct {
	vectors [3]*Vector
	normal  *Vector
	visible bool
	color   uint32
}

func NewTriangle(v1, v2, v3 *Vector, color uint32) *Triangle {
	return &Triangle{
		vectors: [3]*Vector{v1, v2, v3},
		color:   color,
		normal:  NewVector(0, 0, 0),
	}
}

func (t *Triangle) Points() []sdl.FPoint {
	return []sdl.FPoint{
		{X: float32(t.vectors[0].X), Y: float32(t.vectors[0].Y)},
		{X: float32(t.vectors[1].X), Y: float32(t.vectors[1].Y)},
		{X: float32(t.vectors[2].X), Y: float32(t.vectors[2].Y)},
		{X: float32(t.vectors[0].X), Y: float32(t.vectors[0].Y)},
	}
}

func (t *Triangle) Vertices() []sdl.Vertex {
	c := t.FaceColor()
	return []sdl.Vertex{
		{Position: sdl.FPoint{X: float32(t.vectors[0].X), Y: float32(t.vectors[0].Y)}, Color: c, TexCoord: sdl.FPoint{}},
		{Position: sdl.FPoint{X: float32(t.vectors[1].X), Y: float32(t.vectors[1].Y)}, Color: c, TexCoord: sdl.FPoint{}},
		{Position: sdl.FPoint{X: float32(t.vectors[2].X), Y: float32(t.vectors[2].Y)}, Color: c, TexCoord: sdl.FPoint{}},
	}

}

func (t *Triangle) FaceColor() sdl.Color {
	return sdl.Color{
		R: uint8(t.color >> 24),
		G: uint8(t.color >> 16),
		B: uint8(t.color >> 8),
		A: uint8(t.color),
	}
}

func (t *Triangle) Depth() float64 {
	return (t.vectors[0].Z + t.vectors[1].Z + t.vectors[2].Z) / 3
}

type Transformations func(*Triangle) *Triangle

func (t *Triangle) process(f1, f2, f3 *Vector) *Triangle {
	if !t.visible {
		return t
	}
	return &Triangle{
		normal:  t.normal,
		visible: t.visible,
		color:   t.color,
		vectors: [3]*Vector{
			f1,
			f2,
			f3,
		},
	}
}

func WorldMatrices(matrices ...*Matrix4X4) Transformations {
	m1 := IdentityMatrix()
	for _, m := range matrices {
		m1 = m1.MultiplyMatrix(m)
	}
	return func(t *Triangle) *Triangle {
		return t.process(
			t.vectors[0].MatrixMultiply(m1),
			t.vectors[1].MatrixMultiply(m1),
			t.vectors[2].MatrixMultiply(m1),
		)
	}
}

func RotateX(a float64) Transformations {
	return func(t *Triangle) *Triangle {
		rotation := RotateXMatrix(a)
		return t.process(
			t.vectors[0].MatrixMultiply(rotation),
			t.vectors[1].MatrixMultiply(rotation),
			t.vectors[2].MatrixMultiply(rotation),
		)
	}
}

func RotateY(a float64) Transformations {
	return func(t *Triangle) *Triangle {
		rotation := RotateYMatrix(a)
		return t.process(
			t.vectors[0].MatrixMultiply(rotation),
			t.vectors[1].MatrixMultiply(rotation),
			t.vectors[2].MatrixMultiply(rotation),
		)
	}
}

func RotateZ(a float64) Transformations {
	return func(t *Triangle) *Triangle {
		rotation := RotateZMatrix(a)
		return t.process(
			t.vectors[0].MatrixMultiply(rotation),
			t.vectors[1].MatrixMultiply(rotation),
			t.vectors[2].MatrixMultiply(rotation),
		)
	}
}

func Translate(x, y, z float64) Transformations {
	return func(t *Triangle) *Triangle {
		rotation := TranslateMatrix(x, y, z)
		return t.process(
			t.vectors[0].MatrixMultiply(rotation),
			t.vectors[1].MatrixMultiply(rotation),
			t.vectors[2].MatrixMultiply(rotation),
		)
	}
}

func Camera(up, camera, lookDir *Vector, yaw float64) Transformations {
	return func(t *Triangle) *Triangle {
		viewMatrix := LookAtMatrix(camera, camera.Add(lookDir.MatrixMultiply(RotateYMatrix(yaw))), up)
		//		lookDir.Map(lookDir.MatrixMultiply(RotateYMatrix(yaw)))
		//		viewMatrix := LookAtMatrix(camera, camera.Add(lookDir), up)
		return t.process(
			t.vectors[0].MatrixMultiply(viewMatrix),
			t.vectors[1].MatrixMultiply(viewMatrix),
			t.vectors[2].MatrixMultiply(viewMatrix),
		)
	}
}

func Normal(camera *Vector) Transformations {
	return func(t *Triangle) *Triangle {
		if !t.visible {
			return t
		}
		t.normal = t.vectors[1].
			Subtract(t.vectors[0]).
			CrossProduct(t.vectors[2].Subtract(t.vectors[0])).
			Normalize()
		t.visible = t.normal.DotProduct(t.vectors[0].Subtract(camera)) > 0
		return t
	}
}

func Project() Transformations {
	return func(t *Triangle) *Triangle {
		return t.process(
			project()(t.vectors[0]),
			project()(t.vectors[1]),
			project()(t.vectors[2]),
		)
	}
}

func Center(x, y float64) Transformations {
	return func(t *Triangle) *Triangle {
		return t.process(
			center(x, y)(t.vectors[0]),
			center(x, y)(t.vectors[1]),
			center(x, y)(t.vectors[2]),
		)
	}
}

func Shade(light *Vector) Transformations {
	return func(t *Triangle) *Triangle {
		dp := max(0.1, t.normal.DotProduct(light.Normalize()))
		a := uint8(t.color)
		t.color >>= 8
		b := uint8(float64(uint8(t.color)) * dp)
		t.color >>= 8
		g := uint8(float64(uint8(t.color)) * dp)
		t.color >>= 8
		r := uint8(float64(uint8(t.color)) * dp)
		t.color = sdl.Color{R: r, G: g, B: b, A: a}.Uint32()
		return t
	}
}
