package pf

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestMap(t *testing.T) {
	// Define a 6x4 image
	var pitch int32 = 6
	pixels := make([]uint32, 4*pitch)

	// Use 1 core for setting all pixels to 42
	Map(1, func(x uint32) uint32 {
		return 42
	}, pixels)

	// Check that all pixels are now 42
	for i, p := range pixels {
		if p != 42 {
			fmt.Fprintf(os.Stderr, "Map fail at pixel %d, y %d. Has value: %d\n", i, int32(i)/pitch, p)
			t.Fail()
		}
	}
}

func TestConcurrentMap(t *testing.T) {
	// Define a 7x19 image
	var pitch int32 = 7
	pixels := make([]uint32, 19*pitch)

	// Use N cores for setting all pixels to 1337
	n := runtime.NumCPU()
	Map(n, func(x uint32) uint32 {
		return 1337
	}, pixels)

	// Check that all pixels are now 1337
	for i, p := range pixels {
		if p != 1337 {
			fmt.Fprintf(os.Stderr, "Map fail at pixel %d, y %d. Has value: %d\n", i, int32(i)/pitch, p)
			t.Fail()
		}
	}
}

func ExampleConcurrency() {
	// Resolution
	const w, h = 320, 200

	pixels := make([]uint32, w*h)

	// Find the number of available CPUs
	n := runtime.NumCPU()

	// Combine two pixel functions
	pfs := Combine(InvertEverything, OnlyBlue)

	// Run the combined pixel functions over all pixels using all available CPUs
	Map(n, pfs, pixels)

	// Retrieve the red, green and blue components of the first pixel
	red := (pixels[0] | 0x00ff0000) >> 0xffff
	green := (pixels[0] | 0x0000ff00) >> 0xff
	blue := (pixels[0] | 0x000000ff)

	fmt.Printf("rgb(%d, %d, %d)\n", red, green, blue)
	// Output: rgb(0, 0, 255)
}
