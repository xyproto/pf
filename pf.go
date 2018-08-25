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

// Divide a slice of pixels into several slices
func Divide(pixels []uint32, n int) [][]uint32 {
	length := len(pixels)

	sliceLen := length / n
	leftover := length % n

	var sliceOfSlices [][]uint32
	for i := 0; i < (length - leftover); i += sliceLen {
		sliceOfSlices = append(sliceOfSlices, pixels[i:i+sliceLen])
	}
	if leftover > 0 {
		sliceOfSlices = append(sliceOfSlices, pixels[length-leftover:length])
	}
	return sliceOfSlices
}

func Map(cores int, f PixelFunction, pixels []uint32) {
	wg := &sync.WaitGroup{}

	// First copy the pixels into several separate slices
	sliceOfSlices := Divide(pixels, cores)

	// Then process the slices individually
	wg.Add(len(sliceOfSlices))
	for _, subPixels := range sliceOfSlices {
		// subPixels is a slice of pixels
		go func(wg *sync.WaitGroup, subPixels []uint32) {
			for i := range subPixels {
				subPixels[i] = f(subPixels[i])
			}
			wg.Done()
		}(wg, subPixels)
	}
	wg.Wait()

	// Then combine the slices into a new and better slice
	newPixels := make([]uint32, len(pixels))
	for _, subPixels := range sliceOfSlices {
		newPixels = append(newPixels, subPixels...)
	}

	// Finally, replace the pixels with the processed pixels
	pixels = newPixels
}

// GlitchyMap map a PixelFunction over every pixel (uint32 ARGB value).
// This function has data race issues.
func GlitchyMap(cores int, f PixelFunction, pixels []uint32) {
	// Map a pixel function over every pixel, concurrently

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
		for i := int32(0); i < iConcurrentlyDone; i += iStep {
			// run a PixelFunction on parts of the pixel buffer
			wg.Add(1)
			go func(wg *sync.WaitGroup, f PixelFunction, pixels []uint32, iStart, iStop int32) {
				for i := iStart; i < iStop; i++ {
					pixels[i] = f(pixels[i])
				}
				wg.Done()
			}(&wg, f, pixels, i, i+iStep)
		}
		iDone = iConcurrentlyDone
	}

	if iDone == iLength {
		// No leftover pixels
		return
	}

	// Apply partialMap to the final leftover pixels
	wg.Add(1)
	go func(wg *sync.WaitGroup, f PixelFunction, pixels []uint32, iStart, iStop int32) {
		for i := iStart; i < iStop; i++ {
			pixels[i] = f(pixels[i])
		}
		wg.Done()
	}(&wg, f, pixels, iDone, iLength)

	wg.Wait()
}
