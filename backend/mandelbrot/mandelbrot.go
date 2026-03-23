package mandelbrot

import (
	"image"
	"image/color"
	"sync"
)

// Params holds the parameters for rendering a Mandelbrot tile.
type Params struct {
	XMin, XMax float64
	YMin, YMax float64
	Width      int
	Height     int
	MaxIter    int
	ColorScheme string
}

// Render generates a Mandelbrot set image for the given parameters.
// Rows are computed in parallel using goroutines.
func Render(p Params) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, p.Width, p.Height))
	palette := buildPalette(p.ColorScheme, p.MaxIter)

	var wg sync.WaitGroup
	wg.Add(p.Height)

	for py := 0; py < p.Height; py++ {
		py := py
		go func() {
			defer wg.Done()
			y := p.YMin + float64(py)/float64(p.Height)*(p.YMax-p.YMin)
			for px := 0; px < p.Width; px++ {
				x := p.XMin + float64(px)/float64(p.Width)*(p.XMax-p.XMin)
				iter := iterate(x, y, p.MaxIter)
				img.SetNRGBA(px, py, palette[iter])
			}
		}()
	}

	wg.Wait()
	return img
}

// iterate returns the number of iterations before escape, or MaxIter if bounded.
// Uses smooth coloring via the escape-time algorithm with period-2 bulb/cardioid check.
func iterate(cx, cy float64, maxIter int) int {
	// Cardioid / period-2 bulb check for early exit
	q := (cx-0.25)*(cx-0.25) + cy*cy
	if q*(q+(cx-0.25)) <= 0.25*cy*cy {
		return maxIter
	}
	if (cx+1)*(cx+1)+cy*cy <= 0.0625 {
		return maxIter
	}

	zx, zy := 0.0, 0.0
	for i := 0; i < maxIter; i++ {
		zx2, zy2 := zx*zx, zy*zy
		if zx2+zy2 > 4.0 {
			return i
		}
		zy = 2*zx*zy + cy
		zx = zx2 - zy2 + cx
	}
	return maxIter
}

// buildPalette creates a slice of maxIter+1 colors.
// Index maxIter represents the "inside" color (black).
func buildPalette(scheme string, maxIter int) []color.NRGBA {
	palette := make([]color.NRGBA, maxIter+1)
	palette[maxIter] = color.NRGBA{0, 0, 0, 255} // inside = black

	var fn func(t float64) color.NRGBA
	switch scheme {
	case "fire":
		fn = firePalette
	case "ocean":
		fn = oceanPalette
	default:
		fn = classicPalette
	}

	for i := 0; i < maxIter; i++ {
		t := float64(i) / float64(maxIter)
		palette[i] = fn(t)
	}
	return palette
}
