package handler

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter() *gin.Engine {
	r := gin.New()
	r.GET("/api/tile", TileHandler())
	return r
}

const fullQuery = "/api/tile?x_min=-2.5&x_max=1.0&y_min=-1.25&y_max=1.25&width=256&height=256&max_iter=256&color_scheme=classic"

func TestTileHandler_AllParams(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fullQuery, nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "image/png" {
		t.Errorf("Content-Type = %q, want image/png", ct)
	}
	_, err := png.Decode(bytes.NewReader(w.Body.Bytes()))
	if err != nil {
		t.Errorf("response is not valid PNG: %v", err)
	}
}

func TestTileHandler_CustomSize(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet,
		"/api/tile?x_min=-1&x_max=0&y_min=-0.5&y_max=0.5&width=32&height=32&max_iter=64&color_scheme=fire",
		nil,
	)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	img, err := png.Decode(bytes.NewReader(w.Body.Bytes()))
	if err != nil {
		t.Fatalf("response is not valid PNG: %v", err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != 32 || bounds.Dy() != 32 {
		t.Errorf("image size = %dx%d, want 32x32", bounds.Dx(), bounds.Dy())
	}
}

func TestTileHandler_AllColorSchemes(t *testing.T) {
	schemes := []string{"classic", "fire", "ocean"}
	r := newTestRouter()

	for _, scheme := range schemes {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet,
			"/api/tile?x_min=-2.5&x_max=1.0&y_min=-1.25&y_max=1.25&width=16&height=16&max_iter=64&color_scheme="+scheme,
			nil,
		)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("scheme=%q: status = %d, want 200", scheme, w.Code)
		}
	}
}

func TestTileHandler_MissingParams(t *testing.T) {
	r := newTestRouter()
	// Required params: x_min, x_max, y_min, y_max, width, height, max_iter, color_scheme
	cases := []string{
		"/api/tile",
		"/api/tile?x_min=-2.5",
		"/api/tile?x_min=-2.5&x_max=1.0&y_min=-1.25&y_max=1.25&width=64&height=64&max_iter=64",
	}
	for _, query := range cases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, query, nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("query=%q: status = %d, want 400", query, w.Code)
		}
	}
}

func TestTileHandler_InvalidParams(t *testing.T) {
	r := newTestRouter()
	cases := []string{
		"/api/tile?x_min=abc&x_max=1.0&y_min=-1.25&y_max=1.25&width=64&height=64&max_iter=64&color_scheme=classic",
		"/api/tile?x_min=-2.5&x_max=1.0&y_min=-1.25&y_max=1.25&width=notanumber&height=64&max_iter=64&color_scheme=classic",
		"/api/tile?x_min=-2.5&x_max=1.0&y_min=-1.25&y_max=1.25&width=64&height=64&max_iter=xyz&color_scheme=classic",
	}
	for _, query := range cases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, query, nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("query=%q: status = %d, want 400", query, w.Code)
		}
	}
}

func TestTileHandler_CacheHeader(t *testing.T) {
	r := newTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fullQuery, nil)
	r.ServeHTTP(w, req)

	if cc := w.Header().Get("Cache-Control"); cc == "" {
		t.Error("Cache-Control header is missing")
	}
}

// requestTile はテスト用にタイルをリクエストして画像を返す。
func requestTile(t *testing.T, r *gin.Engine, query string) image.Image {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, query, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (query: %s)", w.Code, query)
	}
	img, err := png.Decode(bytes.NewReader(w.Body.Bytes()))
	if err != nil {
		t.Fatalf("PNG デコード失敗: %v", err)
	}
	return img
}

// TestTileHandler_InsideSetIsBlack は集合内（原点付近）が黒であることを確認する。
// 原点 (0+0i) はマンデルブロ集合内のため、対応するピクセルは黒 (0,0,0) になる。
func TestTileHandler_InsideSetIsBlack(t *testing.T) {
	r := newTestRouter()
	// 原点が画像中央に来るよう座標を設定
	// width=101, height=101 として中央ピクセルが (cx=0, cy=0) に対応する
	img := requestTile(t, r,
		"/api/tile?x_min=-1&x_max=1&y_min=-1&y_max=1&width=101&height=101&max_iter=256&color_scheme=classic",
	)
	cx, cy := 50, 50
	r8, g8, b8, _ := img.At(cx, cy).RGBA()
	if r8 != 0 || g8 != 0 || b8 != 0 {
		t.Errorf("原点ピクセル(%d,%d) = (%d,%d,%d), 集合内なので黒を期待", cx, cy, r8>>8, g8>>8, b8>>8)
	}
}

// TestTileHandler_OutsideSetIsColored は集合外（右上遠方）に色が付くことを確認する。
// (2+2i) はマンデルブロ集合の外側のため、対応するピクセルは黒以外になる。
func TestTileHandler_OutsideSetIsColored(t *testing.T) {
	r := newTestRouter()
	// 画像全体が集合外の領域をリクエスト
	img := requestTile(t, r,
		"/api/tile?x_min=1.5&x_max=2.5&y_min=1.5&y_max=2.5&width=16&height=16&max_iter=256&color_scheme=classic",
	)
	// この領域は全ピクセルが脱出するはずなので、少なくとも1ピクセルは非黒であることを確認
	hasColor := false
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r8, g8, b8, _ := img.At(x, y).RGBA()
			if r8 != 0 || g8 != 0 || b8 != 0 {
				hasColor = true
				break
			}
		}
		if hasColor {
			break
		}
	}
	if !hasColor {
		t.Error("集合外の領域なのに全ピクセルが黒")
	}
}

// TestTileHandler_DifferentRegionsDifferentImages は異なる座標範囲で異なる画像が生成されることを確認する。
func TestTileHandler_DifferentRegionsDifferentImages(t *testing.T) {
	r := newTestRouter()
	img1 := requestTile(t, r,
		"/api/tile?x_min=-2.5&x_max=1.0&y_min=-1.25&y_max=1.25&width=32&height=32&max_iter=64&color_scheme=classic",
	)
	img2 := requestTile(t, r,
		"/api/tile?x_min=0.0&x_max=0.5&y_min=0.0&y_max=0.5&width=32&height=32&max_iter=64&color_scheme=classic",
	)

	// 少なくとも1ピクセルが異なることを確認
	different := false
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 {
				different = true
				break
			}
		}
		if different {
			break
		}
	}
	if !different {
		t.Error("異なる座標範囲なのに同一の画像が生成された")
	}
}
