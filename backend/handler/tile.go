package handler

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/osawa-koki/mandelbrot-map/backend/mandelbrot"
)

// TileHandler はマンデルブロ集合タイルを PNG で返す Gin ハンドラーを返す。
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
	xMin, err := requireFloat(c, "x_min")
	if err != nil {
		return mandelbrot.Params{}, err
	}
	xMax, err := requireFloat(c, "x_max")
	if err != nil {
		return mandelbrot.Params{}, err
	}
	yMin, err := requireFloat(c, "y_min")
	if err != nil {
		return mandelbrot.Params{}, err
	}
	yMax, err := requireFloat(c, "y_max")
	if err != nil {
		return mandelbrot.Params{}, err
	}
	width, err := requireInt(c, "width")
	if err != nil {
		return mandelbrot.Params{}, err
	}
	height, err := requireInt(c, "height")
	if err != nil {
		return mandelbrot.Params{}, err
	}
	maxIter, err := requireInt(c, "max_iter")
	if err != nil {
		return mandelbrot.Params{}, err
	}

	scheme, err := requireString(c, "color_scheme")
	if err != nil {
		return mandelbrot.Params{}, err
	}

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

func requireFloat(c *gin.Context, key string) (float64, error) {
	s := c.Query(key)
	if s == "" {
		return 0, fmt.Errorf("missing required parameter: %s", key)
	}
	return strconv.ParseFloat(s, 64)
}

func requireInt(c *gin.Context, key string) (int, error) {
	s := c.Query(key)
	if s == "" {
		return 0, fmt.Errorf("missing required parameter: %s", key)
	}
	return strconv.Atoi(s)
}

func requireString(c *gin.Context, key string) (string, error) {
	s := c.Query(key)
	if s == "" {
		return "", fmt.Errorf("missing required parameter: %s", key)
	}
	return s, nil
}
