package pf

import (
	"sync"
)

// Perform an operation on a single ARGB pixel
type PixelFunction func(v uint32) uint32

// Combine two functions to a single PixelFunction.
// The functions are applied in the same order as the arguments.
func Combine(a, b PixelFunction) PixelFunction {
	return func(v uint32) uint32 {
		return b(a(v))
	}
}

// Combine three functions to a single PixelFunction.
// The functions are applied in the same order as the arguments.
func Combine3(a, b, c PixelFunction) PixelFunction {
	return Combine(a, Combine(b, c))
}

// partialMap runs a PixelFunction on parts of the pixel buffer
func partialMap(wg *sync.WaitGroup, f PixelFunction, sliceOfPixels []uint32, outputPixels chan uint32) {
	defer wg.Done()
	for i := range sliceOfPixels {
		outputPixels <- f(sliceOfPixels[i])
	}
}

// Map a PixelFunction over every pixel (uint32 ARGB value), concurrently
func Map(cores int, f PixelFunction, pixels []uint32) {
	var (
		wg      sync.WaitGroup
		iLength = int32(len(pixels))
		iStep   = iLength / int32(cores)

		// iConcurrentlyDone keeps track of how much work have been done by launching goroutines
		iConcurrentlyDone = int32(cores) * iStep

		// iDone keeps track of how much work have been done in total
		iDone int32
	)

	// Apply partialMap for each of the partitions
	if iStep < iLength {
		var i int32
		for ; i < iConcurrentlyDone; i += iStep {
			wg.Add(1)
			pixelTube := make(chan uint32, iStep)
			go partialMap(&wg, f, pixels[i:i+iStep], pixelTube)
			for n := i; n < i+iStep; n++ {
				pixels[n] = <-pixelTube
			}
		}
		iDone = i
	}

	if iDone == iLength {
		// No leftover pixels
		return
	}

	// Apply partialMap to the final leftover pixels
	wg.Add(1)
	pixelTube := make(chan uint32, iStep)
	go partialMap(&wg, f, pixels[iDone:iLength], pixelTube)
	for n := iDone; n < iLength; n++ {
		pixels[n] = <-pixelTube
	}

	wg.Wait()
}
