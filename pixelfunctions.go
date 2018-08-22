package pf

// The uint32 is on the form ARGB

// Invert the colors
func Invert(v uint32) uint32 {
	// Invert the colors, but set the alpha value to 0xff
	return (0xffffffff - v) | 0xff000000
}

// Keep the red component
func Red(v uint32) uint32 {
	// Keep alpha and the red value
	return v & 0xffff0000
}

// Keep the green component
func Green(v uint32) uint32 {
	// Keep alpha and the green value
	return v & 0xff00ff00
}

// Keep the blue component
func Blue(v uint32) uint32 {
	// Keep alpha and the blue value
	return v & 0xff0000ff
}
