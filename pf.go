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
func (ps *PixelSlice) Map(f PixelFunction) {
	defer ps.wg.Done()
	for i := range ps.pixels {
		ps.pixels = f(ps.pixels[i])
	}
}

type SubPixels struct {
	pixels []uint32
	wg *sync.WaitGroup
	// The start and stop index in the larger context
	iStart uint32
	iStop uint32
}

// Map a PixelFunction over every pixel (uint32 ARGB value), concurrently
func Map(cores int, f PixelFunction, pixels []uint32) {

	var (
		iLength = int32(len(pixels))
		iStep   = iLength / int32(cores)

		// iConcurrentlyDone keeps track of how much work have been done by launching goroutines
		iConcurrentlyDone = int32(cores) * iStep

		// iDone keeps track of how much work have been done in total
		iDone int32
	)

	// Apply partialMap for each of the partitions
	if iStep < iLength {
		wg.Add(cores)
		for i := int32(0); i < iConcurrentlyDone; i += iStep {
			go partialMap(&wg, mut, f, pixels[i:i+iStep])
		}
		iDone = iConcurrentlyDone
	}

	if iDone == iLength {
		// No leftover pixels
		return
	}

	// Apply partialMap to the final leftover pixels
	wg.Add(1)
	go partialMap(&wg, mut, f, pixels[iDone:iLength])

	wg.Wait()
}
