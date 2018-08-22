package pf

import (
	"sync"
)

// Perform an operation on a single ARGB pixel
type PixelFunction func(v uint32) uint32

// Combine two functions to a single PixelFunction
func Combine(a, b PixelFunction) PixelFunction {
	return func(v uint32) uint32 {
		return b(a(v))
	}
}

// partialMap runs a PixelFunction on parts of the pixel buffer
func partialMap(f PixelFunction, pixels []uint32, iStart, iStop int32) {
	for i := iStart; i < iStop; i++ {
		pixels[i] = f(pixels[i])
	}
}

// Map a PixelFunction over every pixel (uint32 ARGB value)
func Map(cores int, f PixelFunction, pixels []uint32) {
	// Map a pixel function over every pixel, concurrently
	var wg sync.WaitGroup

	iLength := int32(len(pixels))

	iStep := iLength / int32(cores)
	iConcurrentlyDone := int32(cores) * iStep

	// Apply partialMap for each of the partitions
	for i := int32(0); i < iConcurrentlyDone; i += iStep {
		wg.Add(1)
		go partialMap(f, pixels, i, i+iStep)
	}

	// Apply partialMap to the final leftover pixels
	wg.Add(1)
	go partialMap(f, pixels, iConcurrentlyDone, iLength)
}
