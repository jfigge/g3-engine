/*
 * Copyright (C) 2023 by Jason Figge
 */

package controller

import (
	"fmt"
	"math"

	"g3-engine/shapes"

	"github.com/jfigge/guilib/graphics"
	"github.com/jfigge/guilib/graphics/fonts"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPSX = 205
	FOV  = 90
	DOV  = 20
)

var (
	Black   = sdl.Color{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	Red     = sdl.Color{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	Yellow  = sdl.Color{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}
	Green   = sdl.Color{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	Cyan    = sdl.Color{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}
	Blue    = sdl.Color{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}
	Magenta = sdl.Color{R: 0xFF, G: 0x00, B: 0xFF, A: 0xFF}
	White   = sdl.Color{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
)

type DirectionCd int

const (
	DirectionCdForward DirectionCd = iota
	DirectionCdBackward
	DirectionCdAntiClockwise
	DirectionCdClockwise
	DirectionCdMoveUp
	DirectionCdMoveDown
	DirectionCdLookUp
	DirectionCdLookDown
	DirectionCdStrafeLeft
	DirectionCdStrafeRight
)

type Camera struct {
	up      *shapes.Vector
	camera  *shapes.Vector
	lookDir *shapes.Vector
	yaw     float64
	light   *shapes.Vector
}

type Fov struct {
	width  float64
	height float64
	cw     float64
	ch     float64
	ndov   float64 // near depth of view
	fdov   float64 // far depth of view
}

type Controller struct {
	graphics.BaseHandler
	graphics.CoreMethods
	camera *Camera
	fov    *Fov
	shapes []*shapes.Shape
}

func NewController(width, height float64) *Controller {
	c := &Controller{
		camera: &Camera{
			up:      shapes.NewVector(0, 1, 0),
			camera:  shapes.NewVector(0, 0, 0),
			lookDir: shapes.NewVector(0, 0, 1),
			light:   shapes.NewVector(0, 0, -1),
			yaw:     -math.Pi,
		},
		fov: &Fov{
			width:  width,
			height: height,
			cw:     width / 2,
			ch:     height / 2,
		},
	}
	f := FOV * math.Pi / 360
	c.fov.ndov = 0.1  //c.fov.cw * math.Tan(f)
	c.fov.fdov = 1000 //c.fov.ndov * DOV
	a := width / height

	shapes := shapes.LoadShapes(shapes.ProjectMatrix(a, 1/f, c.fov.ndov, c.fov.fdov))
	c.shapes = append(c.shapes,
		shapes.Axis(),
		//Locate(c.fov.cw, c.fov.ch, c.fov.ndov*2).
		//Rotate(22, 44, 66).
		//Scale(50, 50, 50),
	)
	return c
}

func (c *Controller) Init(canvas *graphics.Canvas) {
	fonts.LoadFonts(canvas.Renderer())
	graphics.ErrorTrap(canvas.Renderer().SetDrawBlendMode(sdl.BLENDMODE_BLEND))
	canvas.Renderer().SetLogicalSize(int32(c.fov.width), int32(c.fov.height))
	c.AddDestroyer(fonts.FreeFonts)
}

func (c *Controller) OnDraw(renderer *sdl.Renderer) {
	graphics.ErrorTrap(c.Clear(renderer, uint32(0x232323)))
	c.draw3D(renderer)
	graphics.ErrorTrap(c.WriteFrameRate(renderer, FPSX, 0))
	fonts.Default.PrintfAt(renderer, 20, 20, 0xFFFFFF, fmt.Sprintf("Camera: %f, %f, %f", c.camera.camera.X, c.camera.camera.Y, c.camera.camera.Z))
}

func (c *Controller) OnUpdate() {
	c.processKeys()
}

var xa float64

func (c *Controller) draw3D(renderer *sdl.Renderer) {
	//xa = xa + .01
	for _, shape := range c.shapes {
		ts := shape.GetTriangles(
			shapes.WorldMatrices(
				shapes.RotateXMatrix(xa*1/2),
				shapes.RotateYMatrix(xa*2/3),
				shapes.RotateZMatrix(xa),
				shapes.TranslateMatrix(0, 0, 9),
			),
			shapes.Normal(c.camera.camera),
			shapes.Camera(c.camera.up, c.camera.camera, c.camera.lookDir, c.camera.yaw),
			shapes.Project(),
			shapes.Center(c.fov.cw, c.fov.ch),
			shapes.Shade(c.camera.light),
		)
		//slices.SortFunc(ts, func(t1 *shapes.Triangle, t2 *shapes.Triangle) int {
		//	return int(t1.Depth() - t2.Depth())
		//})

		var vs []sdl.Vertex
		for _, t := range ts {
			vs = append(vs, t.Vertices()...)
		}

		renderer.RenderGeometry(nil, vs, nil)

		//renderer.SetDrawColor(0, 0, 0, 0xFF)
		//for _, t := range ts {
		//	pts := t.GetPoints()
		//	renderer.DrawLinesF(pts)
		//}
	}
}

func (c *Controller) processKeys() {
	codes := sdl.GetKeyboardState()
	//shift := codes[sdl.SCANCODE_LSHIFT] == 1 || codes[sdl.SCANCODE_RSHIFT] == 1
	if codes[sdl.SCANCODE_UP] == 1 {
		c.move(DirectionCdMoveUp)
	} else if codes[sdl.SCANCODE_DOWN] == 1 {
		c.move(DirectionCdMoveDown)
	}
	if codes[sdl.SCANCODE_LEFT] == 1 {
		c.move(DirectionCdStrafeLeft)
	} else if codes[sdl.SCANCODE_RIGHT] == 1 {
		c.move(DirectionCdStrafeRight)
	}
	if codes[sdl.SCANCODE_W] == 1 {
		c.move(DirectionCdForward)
	} else if codes[sdl.SCANCODE_S] == 1 {
		c.move(DirectionCdBackward)
	}
	if codes[sdl.SCANCODE_A] == 1 {
		c.move(DirectionCdAntiClockwise)
	} else if codes[sdl.SCANCODE_D] == 1 {
		c.move(DirectionCdClockwise)
	}
}

func (c *Controller) move(dir DirectionCd) {
	switch dir {
	case DirectionCdForward:
		c.camera.camera.Map(c.camera.camera.Add(c.camera.lookDir.Multiply(.2)))
	case DirectionCdBackward:
		c.camera.camera.Map(c.camera.camera.Subtract(c.camera.lookDir.Multiply(.2)))
	case DirectionCdStrafeLeft:
		c.camera.camera.X += .2
	case DirectionCdStrafeRight:
		c.camera.camera.X -= .2
	case DirectionCdMoveUp:
		c.camera.camera.Y += .2
	case DirectionCdMoveDown:
		c.camera.camera.Y -= .2
	case DirectionCdLookUp:
		//c.camera.l += 1
		//if c.camera.l > cameraHeight {
		//	c.camera.l = cameraHeight
		//}
	case DirectionCdLookDown:
		//c.camera.l -= 1
		//if c.camera.l < -cameraHeight {
		//	c.camera.l = -cameraHeight
		//}
	case DirectionCdAntiClockwise:
		c.camera.yaw -= .01
	case DirectionCdClockwise:
		c.camera.yaw += .01

	}
}
