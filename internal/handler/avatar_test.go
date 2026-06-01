package handler

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"
)

// multipartImage builds an in-memory multipart/form-data body with a single
// "file" field containing the given bytes, returning the body and content type.
func multipartImage(t *testing.T, content []byte) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fw, err := mw.CreateFormFile("file", "upload.bin")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := mw.Close(); err != nil {
		t.Fatal(err)
	}
	return body, mw.FormDataContentType()
}

func TestReadUploadedImageAcceptsPNG(t *testing.T) {
	dir := t.TempDir()
	png := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0, 0, 0, 0}
	body, contentType := multipartImage(t, png)

	r := httptest.NewRequest("POST", "/characters/1/avatar", body)
	r.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	rel, ok := readUploadedImage(w, r, dir, "character")
	if !ok {
		t.Fatalf("esperaba ok=true; status=%d body=%s", w.Code, w.Body.String())
	}
	if !strings.HasSuffix(rel, ".png") {
		t.Fatalf("esperaba un .png, obtuve %q", rel)
	}
}

func TestReadUploadedImageRejectsNonImage(t *testing.T) {
	dir := t.TempDir()
	body, contentType := multipartImage(t, []byte("solo texto, no una imagen real"))

	r := httptest.NewRequest("POST", "/characters/1/avatar", body)
	r.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	if _, ok := readUploadedImage(w, r, dir, "character"); ok {
		t.Fatal("esperaba que rechazara un archivo que no es imagen")
	}
	if w.Code != 400 {
		t.Fatalf("esperaba 400, obtuve %d", w.Code)
	}
}
