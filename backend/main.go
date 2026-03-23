package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/osawa-koki/mandelbrot-map/backend/handler"
)

func main() {
	r := newRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func newRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/api/tile", handler.TileHandler())

	// フロントエンドの静的ファイルを配信する（./dist が存在する場合のみ）
	if _, err := os.Stat("./dist"); err == nil {
		r.Static("/assets", "./dist/assets")
		r.NoRoute(func(c *gin.Context) {
			c.File("./dist/index.html")
		})
	}

	return r
}
