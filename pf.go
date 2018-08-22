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
func partialMap(f PixelFunction, pixels []uint32, pitch, startY, stopY int32) {
	startPos := startY * pitch
	stopPos := stopY*pitch + pitch
	for i := startPos; i < stopPos; i++ {
		pixels[i] = f(pixels[i])
	}
}

// Map a PixelFunction over every pixel (uint32 ARGB value)
func Map(cores int, f PixelFunction, pixels []uint32, pitch int32) {
	// Map a pixel function over every pixel, concurrently
	var wg sync.WaitGroup

	height := int32(len(pixels)) / pitch
	ystep := height / int32(cores)
	ymax := ystep * int32(cores)

	// Apply partialMap for each of the partitions along the y axis
	for y := int32(0); y < ymax; y += ystep {
		wg.Add(1)
		go partialMap(f, pixels, pitch, y, y+ystep)
	}

	// Apply partialMap to the final leftover pixels
	wg.Add(1)
	go partialMap(f, pixels, pitch, ymax, height)
}
