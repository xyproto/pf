# pf - PixelFunction

## Work in progress, needs more testing!

# Description

The PixelFunction type has this signature:

    func(v uint32) uint32

If you have a pixel buffer of type []uint32, with colors on the form ARGB, then this modules allows you to apply PixelFunctions to that slice, concurrently.

PixelFunctions can also be combined to a single PixelFunction.

The goal is to avoid looping over all pixels more than once, while applying many different effects, concurrently.
