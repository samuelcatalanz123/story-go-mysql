// Package upload validates and stores uploaded image files on local disk.
// It never trusts the client-provided Content-Type or filename: the MIME type
// is sniffed from the file's first bytes, and the stored name is randomly
// generated to avoid path-traversal attacks.
package upload

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// MaxImageBytes is the largest image we accept (5 MB).
const MaxImageBytes = 5 << 20

// ErrUnsupportedType is returned when the file is not a JPEG, PNG or WebP.
var ErrUnsupportedType = errors.New("unsupported image type")

// allowedTypes maps the sniffed MIME type to the file extension we store.
var allowedTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// SaveImage validates the image in r and writes it under baseDir/subdir with a
// random name. It returns the path relative to baseDir (e.g. "avatars/ab12.jpg")
// so the caller can build the public URL. The MIME type is detected from the
// first 512 bytes, not from any client header.
func SaveImage(baseDir, subdir string, r io.Reader) (string, error) {
	// Sniff the content type from the first 512 bytes (what DetectContentType
	// needs), then prepend them back so the whole file still gets written.
	head := make([]byte, 512)
	n, err := io.ReadFull(r, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return "", err
	}
	head = head[:n]

	ext, ok := allowedTypes[http.DetectContentType(head)]
	if !ok {
		return "", ErrUnsupportedType
	}

	name, err := randomName()
	if err != nil {
		return "", err
	}
	filename := name + ext

	dir := filepath.Join(baseDir, subdir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	dst, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Stream the file to disk: the sniffed head first, then the rest. io.Copy
	// streams in chunks instead of loading the whole file into memory.
	if _, err := io.Copy(dst, io.MultiReader(bytes.NewReader(head), r)); err != nil {
		return "", err
	}

	// Forward slash so the value works directly inside a URL path.
	return subdir + "/" + filename, nil
}

// randomName returns a 32-character random hex string, safe to use as a
// filename (no client input, so no path traversal possible).
func randomName() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
