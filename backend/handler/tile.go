package handler

import (
	"bytes"
	"image/png"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/osawa-koki/mandelbrot-map/backend/mandelbrot"
)

const (
	defaultMaxIter    = 256
	defaultWidth      = 256
	defaultHeight     = 256
	defaultColorScheme = "classic"
)

// TileHandler returns a Gin handler that renders a Mandelbrot tile as PNG.
func TileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, err := parseTileParams(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		img := mandelbrot.Render(p)

		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode image"})
			return
		}

		c.Header("Cache-Control", "public, max-age=86400")
		c.Data(http.StatusOK, "image/png", buf.Bytes())
	}
}

func parseTileParams(c *gin.Context) (mandelbrot.Params, error) {
	xMin, err := parseFloat(c, "x_min", -2.5)
	if err != nil {
		return mandelbrot.Params{}, err
	}
	xMax, err := parseFloat(c, "x_max", 1.0)
	if err != nil {
		return mandelbrot.Params{}, err
	}
	yMin, err := parseFloat(c, "y_min", -1.25)
	if err != nil {
		return mandelbrot.Params{}, err
	}
	yMax, err := parseFloat(c, "y_max", 1.25)
	if err != nil {
		return mandelbrot.Params{}, err
	}
	width, err := parseInt(c, "width", defaultWidth)
	if err != nil {
		return mandelbrot.Params{}, err
	}
	height, err := parseInt(c, "height", defaultHeight)
	if err != nil {
		return mandelbrot.Params{}, err
	}
	maxIter, err := parseInt(c, "max_iter", defaultMaxIter)
	if err != nil {
		return mandelbrot.Params{}, err
	}

	scheme := c.DefaultQuery("color_scheme", defaultColorScheme)

	return mandelbrot.Params{
		XMin:        xMin,
		XMax:        xMax,
		YMin:        yMin,
		YMax:        yMax,
		Width:       width,
		Height:      height,
		MaxIter:     maxIter,
		ColorScheme: scheme,
	}, nil
}

func parseFloat(c *gin.Context, key string, def float64) (float64, error) {
	s := c.Query(key)
	if s == "" {
		return def, nil
	}
	return strconv.ParseFloat(s, 64)
}

func parseInt(c *gin.Context, key string, def int) (int, error) {
	s := c.Query(key)
	if s == "" {
		return def, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}
