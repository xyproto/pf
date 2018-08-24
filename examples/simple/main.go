package main

import (
	"fmt"
	"github.com/xyproto/pf"
	"runtime"
)

func main() {
	// Resolution
	const w, h = 320, 200

	pixels := make([]uint32, w*h)

	// Find the number of available CPUs
	n := runtime.NumCPU()

	// Combine two pixel functions
	pfs := pf.Combine(pf.InvertEverything, pf.OnlyGreen)

	// Run the combined pixel functions over all pixels using all available CPUs
	pf.Map(n, pfs, pixels)

	// Retrieve the red, green and blue components of the first pixel
	red := (pixels[0] & 0x00ff0000) >> 16
	green := (pixels[0] & 0x0000ff00) >> 8
	blue := (pixels[0] & 0x000000ff)

	// Should output only blue: rgb(0, 255, 0)
	fmt.Printf("rgb(%d, %d, %d) / #%x\n", red, green, blue, pixels[0])
}
