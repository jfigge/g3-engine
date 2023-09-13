/*
 * Copyright (C) 2023 by Jason Figge
 */

package main

import (
	"fmt"

	"g3-engine/code/controller"

	"github.com/jfigge/guilib/graphics"
)

const (
	screenWidth  = 720
	screenHeight = 360
	Scale        = 2
)

func main() {
	graphics.Open(
		"g3 engine",
		screenWidth*Scale,
		screenHeight*Scale,
		controller.NewController(screenWidth/2, screenHeight),
		graphics.Framerate(60),
	)
	fmt.Println("Game over")
}
