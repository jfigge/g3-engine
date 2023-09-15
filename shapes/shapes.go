/*
 * Copyright (C) 2023 by Jason Figge
 */

package shapes

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

var projectionMatrix *Matrix4X4

type Shapes struct {
	shapes map[string]*Shape
}

func LoadShapes(pm *Matrix4X4) *Shapes {
	projectionMatrix = pm
	s := &Shapes{
		shapes: map[string]*Shape{},
	}
	s.shapes["cube"] = createCube()
	s.shapes["spaceship"] = loadObject("spaceship.obj")
	s.shapes["teapot"] = loadObject("teapot.obj")
	s.shapes["axis"] = loadObject("axis.obj")
	return s
}

func (s *Shapes) Cube() *Shape {
	return s.shapes["cube"].duplicate()
}

func (s *Shapes) Spaceship() *Shape {
	return s.shapes["spaceship"].duplicate()
}

func (s *Shapes) Teapot() *Shape {
	return s.shapes["teapot"].duplicate()
}

func (s *Shapes) Axis() *Shape {
	return s.shapes["axis"].duplicate()
}

func createCube() *Shape {
	pts := []*Vector{
		//{-1, 1, 1},
		//{-1, -1, 1},
		//{1, -1, 1},
		//{1, 1, 1},
		//{-1, 1, -1},
		//{-1, -1, -1},
		//{1, -1, -1},
		//{1, 1, -1},
		{0, 0, 0, 1},
		{0, 1, 0, 1},
		{1, 1, 0, 1},
		{1, 0, 0, 1},
		{0, 0, 1, 1},
		{0, 1, 1, 1},
		{1, 1, 1, 1},
		{1, 0, 1, 1},
	}
	idx := []int{
		// Front
		0, 1, 2, 0, 0, 2, 3, 0,
		// Left
		4, 5, 1, 4, 4, 1, 0, 4,
		// Back
		7, 6, 5, 7, 7, 5, 4, 7,
		// Right
		3, 2, 6, 3, 3, 6, 7, 3,
		// Top
		1, 5, 6, 1, 1, 6, 2, 1,
		//Bottom
		4, 0, 3, 4, 4, 3, 7, 4,
	}
	clr := []uint32{
		0xFF0000FF,
		0xFFFF00FF,
		0x00FF00FF,
		0x00FFFFFF,
		0x0000FFFF,
		0xFF00FFFF,
	}
	s := &Shape{
		ts:       make([]*Triangle, len(idx)/4),
		location: NewVector(0, 0, 0),
		rotation: NewVector(0, 0, 0),
		scale:    NewVector(1, 1, 1),
	}
	for i := 0; i < len(idx)/4; i++ {
		s.ts[i] = NewTriangle(
			pts[idx[i*4+0]],
			pts[idx[i*4+1]],
			pts[idx[i*4+2]],
			clr[i/2],
		)
	}

	return s
}

func loadObject(filename string) *Shape {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("unable to get working directory: %w", err))
	}

	var file *os.File
	base := filepath.Join(dir, "resources", "objects")
	file, err = os.Open(filepath.Join(base, filename))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var pts []*Vector
	var ts []*Triangle

	lineCnt := 0
	for scanner.Scan() {
		lineCnt++
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		switch line[:2] {
		case "v ":
			pts = append(pts, parseVector(line[2:], lineCnt))
		case "f ":
			ts = append(ts, parseFace(pts, line[2:], lineCnt))
		}
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &Shape{
		ts:       ts,
		location: NewVector(0, 0, 0),
		rotation: NewVector(0, 0, 0),
		scale:    NewVector(1, 1, 1),
		color:    sdl.Color{R: uint8(0xff), G: uint8(0xff), B: uint8(0xff), A: uint8(0xff)},
	}
}

func parseVector(line string, lineCnt int) *Vector {
	xyz := strings.Split(line, " ")
	if len(xyz) != 3 && len(xyz) != 4 {
		log.Panicf("Bad vector in line %d: %s", lineCnt, line)
	}
	var x, y, z float64
	var err error

	x, err = strconv.ParseFloat(xyz[0], 32)
	if err != nil {
		log.Panicf("Bad x in line %d: %s", lineCnt, line)
	}

	y, err = strconv.ParseFloat(xyz[1], 32)
	if err != nil {
		log.Panicf("Bad Y in line %d: %s", lineCnt, line)
	}

	z, err = strconv.ParseFloat(xyz[2], 32)
	if err != nil {
		log.Panicf("Bad Z in line %d: %s", lineCnt, line)
	}

	return NewVector(x, y, z)
}

func parseFace(pts []*Vector, line string, lineCnt int) *Triangle {
	xyz := strings.Split(line, " ")
	if len(xyz) != 3 {
		log.Panicf("Bad face in line %d: %s", lineCnt, line)
	}
	var x, y, z int64
	var err error

	x, err = strconv.ParseInt(xyz[0], 10, 32)
	if err != nil {
		log.Panicf("Bad x in line %d: %s", lineCnt, line)
	}

	y, err = strconv.ParseInt(xyz[1], 10, 32)
	if err != nil {
		log.Panicf("Bad Y in line %d: %s", lineCnt, line)
	}

	z, err = strconv.ParseInt(xyz[2], 10, 32)
	if err != nil {
		log.Panicf("Bad Z in line %d: %s", lineCnt, line)
	}

	return &Triangle{
		vectors: [3]*Vector{pts[x-1], pts[y-1], pts[z-1]},
		normal:  NewVector(0, 0, 0),
		visible: true,
		color:   uint32(0xFFFFFFFF),
	}
}
