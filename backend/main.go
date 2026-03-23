package main

import (
	"github.com/gin-gonic/gin"
	"github.com/osawa-koki/mandelbrot-map/backend/handler"
)

func main() {
	r := newRouter()
	r.Run(":8080")
}

func newRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/api/tile", handler.TileHandler())

	return r
}
