package mandelbrot

import (
	"testing"
)

func TestIterate_InsideSet(t *testing.T) {
	// Origin (0,0) should never escape
	if got := iterate(0, 0, 256); got != 256 {
		t.Errorf("iterate(0,0,256) = %d, want 256", got)
	}
	// (−1, 0) is inside the set
	if got := iterate(-1, 0, 256); got != 256 {
		t.Errorf("iterate(-1,0,256) = %d, want 256", got)
	}
}

func TestIterate_OutsideSet(t *testing.T) {
	// (2, 2) escapes immediately
	iter := iterate(2, 2, 256)
	if iter >= 256 {
		t.Errorf("iterate(2,2,256) = %d, expected < 256", iter)
	}
}

func TestIterate_CardioidBulb(t *testing.T) {
	// Interior of main cardioid should reach maxIter
	if got := iterate(0.25, 0, 256); got != 256 {
		t.Errorf("iterate(0.25,0,256) = %d, want 256", got)
	}
	// Period-2 bulb center
	if got := iterate(-1, 0, 256); got != 256 {
		t.Errorf("iterate(-1,0,256) = %d, want 256", got)
	}
}

func TestRender_Size(t *testing.T) {
	p := Params{
		XMin: -2.5, XMax: 1.0,
		YMin: -1.25, YMax: 1.25,
		Width: 64, Height: 48,
		MaxIter:     64,
		ColorScheme: "classic",
	}
	img := Render(p)
	bounds := img.Bounds()
	if bounds.Dx() != 64 || bounds.Dy() != 48 {
		t.Errorf("image size = %dx%d, want 64x48", bounds.Dx(), bounds.Dy())
	}
}

func TestRender_AllColorSchemes(t *testing.T) {
	schemes := []string{"classic", "fire", "ocean", "unknown"}
	for _, scheme := range schemes {
		p := Params{
			XMin: -2.5, XMax: 1.0,
			YMin: -1.25, YMax: 1.25,
			Width: 16, Height: 16,
			MaxIter:     32,
			ColorScheme: scheme,
		}
		img := Render(p)
		if img == nil {
			t.Errorf("Render with scheme=%q returned nil", scheme)
		}
	}
}

func TestBuildPalette_Length(t *testing.T) {
	for _, maxIter := range []int{10, 100, 256} {
		p := buildPalette("classic", maxIter)
		if len(p) != maxIter+1 {
			t.Errorf("buildPalette len = %d, want %d", len(p), maxIter+1)
		}
		// Inside color must be black
		inside := p[maxIter]
		if inside.R != 0 || inside.G != 0 || inside.B != 0 || inside.A != 255 {
			t.Errorf("inside color = %v, want black", inside)
		}
	}
}

func BenchmarkRender(b *testing.B) {
	p := Params{
		XMin: -2.5, XMax: 1.0,
		YMin: -1.25, YMax: 1.25,
		Width: 256, Height: 256,
		MaxIter:     256,
		ColorScheme: "classic",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Render(p)
	}
}
