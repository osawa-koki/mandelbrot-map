package mandelbrot

import (
	"image/color"
	"math"
)

// classicPalette maps t in [0,1) to a blue-white-black gradient.
func classicPalette(t float64) color.NRGBA {
	// Smooth band coloring
	r := uint8(9 * (1 - t) * t * t * t * 255)
	g := uint8(15 * (1 - t) * (1 - t) * t * t * 255)
	b := uint8(8.5 * (1 - t) * (1 - t) * (1 - t) * t * 255)
	return color.NRGBA{r, g, b, 255}
}

// firePalette maps t in [0,1) to a black-red-orange-yellow-white gradient.
func firePalette(t float64) color.NRGBA {
	r := clampUint8(math.Pow(t, 0.5) * 255 * 2)
	g := clampUint8(math.Pow(t, 1.5) * 255 * 2)
	b := clampUint8(math.Pow(t, 4.0) * 255 * 2)
	return color.NRGBA{r, g, b, 255}
}

// oceanPalette maps t in [0,1) to a deep-blue to cyan gradient.
func oceanPalette(t float64) color.NRGBA {
	r := clampUint8(t * 30)
	g := clampUint8(t * 180)
	b := clampUint8(50 + t*205)
	return color.NRGBA{r, g, b, 255}
}

func clampUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
