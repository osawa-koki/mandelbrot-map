package mandelbrot

import (
	"image/color"
	"math"
)

// classicPalette は t∈[0,1) を青〜白〜黒のグラデーションにマップする。
func classicPalette(t float64) color.NRGBA {
	// スムーズバンドカラーリング
	r := uint8(9 * (1 - t) * t * t * t * 255)
	g := uint8(15 * (1 - t) * (1 - t) * t * t * 255)
	b := uint8(8.5 * (1 - t) * (1 - t) * (1 - t) * t * 255)
	return color.NRGBA{r, g, b, 255}
}

// firePalette は t∈[0,1) を黒→赤→オレンジ→黄→白のグラデーションにマップする。
func firePalette(t float64) color.NRGBA {
	r := clampUint8(math.Pow(t, 0.5) * 255 * 2)
	g := clampUint8(math.Pow(t, 1.5) * 255 * 2)
	b := clampUint8(math.Pow(t, 4.0) * 255 * 2)
	return color.NRGBA{r, g, b, 255}
}

// oceanPalette は t∈[0,1) を深海ブルー〜シアンのグラデーションにマップする。
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
