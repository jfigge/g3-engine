/*
 * Copyright (C) 2023 by Jason Figge
 */

package shapes

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Shape struct {
	ts       []*Triangle
	location *Vector
	rotation *Vector
	scale    *Vector
	color    sdl.Color
}

func (s *Shape) duplicate() *Shape {
	s2 := &Shape{
		ts:       make([]*Triangle, len(s.ts)),
		location: NewVector(s.location.X, s.location.Y, s.location.Z),
		rotation: NewVector(s.rotation.X, s.rotation.Y, s.rotation.Z),
		scale:    NewVector(s.scale.X, s.scale.Y, s.scale.Z),
	}
	for i, t := range s.ts {
		s2.ts[i] = NewTriangle(t.vectors[0], t.vectors[1], t.vectors[2], t.color)
	}
	return s2
}

func (s *Shape) Locate(x, y, z float64) *Shape {
	s.location = NewVector(x, y, z)
	return s
}

func (s *Shape) Rotate(x, y, z float64) *Shape {
	s.rotation = NewVector(x, y, z)
	return s
}

func (s *Shape) Scale(x, y, z float64) *Shape {
	s.scale = NewVector(x, y, z)
	return s
}

func (s *Shape) GetTriangles(transforms ...Transformations) []*Triangle {
	ts := make([]*Triangle, 0, len(s.ts))
	for _, t := range s.ts {
		t.visible = true
		t2 := t
		for _, transform := range transforms {
			t2 = transform(t2)
		}
		if t2.visible {
			ts = append(ts, t2)
		}
	}
	return ts
}
