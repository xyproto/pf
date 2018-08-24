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

func TestDivide(t *testing.T) {
	pixels := make([]uint32, 100)

	// 100 divided into 2 slices fits perfectly
	slicesOfSlices := Divide(pixels, 2)
	if len(slicesOfSlices) != 2 {
		t.Fail()
	}
	if len(slicesOfSlices[0]) != 50 {
		t.Fail()
	}
	if len(slicesOfSlices[1]) != 50 {
		t.Fail()
	}

	// 100 divided into 3 slices + leftover pixels
	slicesOfSlices = Divide(pixels, 3)
	if len(slicesOfSlices) != 4 {
		t.Fail()
	}
	if len(slicesOfSlices[0]) != 33 {
		t.Fail()
	}
	if len(slicesOfSlices[1]) != 33 {
		t.Fail()
	}
	if len(slicesOfSlices[2]) != 33 {
		t.Fail()
	}
	if len(slicesOfSlices[3]) != 1 {
		t.Fail()
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
