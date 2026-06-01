package upload

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// pngHeader is a minimal valid PNG signature so http.DetectContentType
// reports "image/png".
var pngHeader = []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0, 0, 0, 0}

func TestSaveImageStoresPNG(t *testing.T) {
	dir := t.TempDir()
	rel, err := SaveImage(dir, "avatars", bytes.NewReader(pngHeader))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(rel, "avatars/") || !strings.HasSuffix(rel, ".png") {
		t.Fatalf("ruta relativa inesperada: %s", rel)
	}
	if _, err := os.Stat(filepath.Join(dir, rel)); err != nil {
		t.Fatalf("el archivo no se guardó: %v", err)
	}
}

func TestSaveImageRejectsNonImage(t *testing.T) {
	dir := t.TempDir()
	_, err := SaveImage(dir, "avatars", strings.NewReader("esto es texto, no una imagen"))
	if !errors.Is(err, ErrUnsupportedType) {
		t.Fatalf("esperaba ErrUnsupportedType, obtuve %v", err)
	}
}
