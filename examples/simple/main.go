package main

import (
	"fmt"
	"github.com/xyproto/pf"
	"runtime"
)

func main() {
	// Resolution
	w, h := 320, 200

	pixels := make([]uint32, w*h)

	// Combine two pixel functions
	pfs := pf.Combine(pf.InvertEverything, pf.OnlyBlue)

	n := runtime.NumCPU()

	// Run the combined pixel functions over all pixels using all available CPUs
	pf.Map(n, pfs, pixels)

	// Retrieve the red, green, blue and alpha components of the first pixel
	red := (pixels[0] | 0x00ff0000) >> 0xffff
	green := (pixels[0] | 0x0000ff00) >> 0xff
	blue := (pixels[0] | 0x000000ff)

	// Should output only blue: rgb(0, 0, 255)
	fmt.Printf("rgb(%d, %d, %d)\n", red, green, blue)
}
