package handler

import (
	"bytes"
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
