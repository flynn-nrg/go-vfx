package oiio

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// writeFlatTestTexture writes a small flat-color PNG so filtered lookups are
// deterministic regardless of the mip level or filter width OIIO picks. The
// alpha must be non-opaque: Go's png encoder silently drops the alpha channel
// (encoding RGB instead of RGBA) whenever the image is fully opaque.
func writeFlatTestTexture(t *testing.T, path string, c color.NRGBA) {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	for y := range 8 {
		for x := range 8 {
			img.SetNRGBA(x, y, c)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create test texture: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("failed to encode test texture: %v", err)
	}
}

func TestTextureSystemLookup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "flat.png")
	writeFlatTestTexture(t, path, color.NRGBA{R: 255, G: 128, B: 0, A: 254})

	ts, err := NewTextureSystem()
	if err != nil {
		t.Fatalf("NewTextureSystem: %v", err)
	}
	defer ts.Close()

	opts := NewTextureLookupOptions()
	result, err := ts.Texture(path, opts, 0.5, 0.5, 0, 0, 0, 0, 4)
	if err != nil {
		t.Fatalf("Texture: %v", err)
	}
	if len(result) != 4 {
		t.Fatalf("expected 4 channels, got %d", len(result))
	}

	// PNG is 8-bit sRGB; OIIO decodes it to linear-ish float unless told
	// otherwise, so just sanity-check ordering (R > G > B) and full alpha
	// rather than exact values.
	if !(result[0] > result[1] && result[1] > result[2]) {
		t.Fatalf("expected R > G > B, got %v", result)
	}
	if result[3] < 0.9 {
		t.Fatalf("expected alpha close to 1.0, got %v", result[3])
	}
}

func TestTextureSystemGetTextureInfo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "flat.png")
	writeFlatTestTexture(t, path, color.NRGBA{R: 10, G: 20, B: 30, A: 254})

	ts, err := NewTextureSystem()
	if err != nil {
		t.Fatalf("NewTextureSystem: %v", err)
	}
	defer ts.Close()

	channels, err := ts.GetTextureInfoInt(path, "channels")
	if err != nil {
		t.Fatalf("GetTextureInfoInt(channels): %v", err)
	}
	if channels != 4 {
		t.Fatalf("expected 4 channels, got %d", channels)
	}

	exists, err := ts.GetTextureInfoInt(path, "exists")
	if err != nil {
		t.Fatalf("GetTextureInfoInt(exists): %v", err)
	}
	if exists != 1 {
		t.Fatalf("expected exists=1, got %d", exists)
	}
}

func TestTextureSystemMissingFile(t *testing.T) {
	ts, err := NewTextureSystem()
	if err != nil {
		t.Fatalf("NewTextureSystem: %v", err)
	}
	defer ts.Close()

	opts := NewTextureLookupOptions()
	_, err = ts.Texture("/nonexistent/path/does-not-exist.png", opts, 0.5, 0.5, 0, 0, 0, 0, 4)
	if err == nil {
		t.Fatalf("expected an error looking up a missing texture")
	}
}
