package pf

import (
	"os"
	"testing"
	"fmt"
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
			fmt.Fprintf(os.Stderr, "Map fail at pixel %d, y %d. Has value: %d\n", i, int32(i) / pitch, p)
			t.Fail()
		}
	}
}
