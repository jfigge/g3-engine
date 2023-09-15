/*
 * Copyright (C) 2023 by Jason Figge
 */

package shapes

import (
	"fmt"
	"testing"
)

func TestDotProduct(t *testing.T) {
	v0 := NewVector(0, 0, 1)
	v1 := NewVector(0, 0, 0.5)
	v2 := NewVector(0, 0, 0.5)
	v3 := NewVector(0, 0, 0.5)

	fmt.Printf("%f\n", v1.DotProduct(v0))
	fmt.Printf("%f\n", v2.DotProduct(v0))
	fmt.Printf("%f\n", v3.DotProduct(v0))
}
