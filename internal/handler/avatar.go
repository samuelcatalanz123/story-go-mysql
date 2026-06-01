package handler

import (
	"errors"
	"net/http"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/upload"
	"story-go-mysql/internal/web"
)

// readUploadedImage reads the "file" field of a multipart/form-data request,
// validates it (size and image type) and stores it under uploadDir/avatars.
// It returns the path relative to uploadDir (e.g. "avatars/ab12.jpg") and ok.
// On any failure it writes the HTTP error response and returns ok=false.
func readUploadedImage(w http.ResponseWriter, r *http.Request, uploadDir, resource string) (string, bool) {
	// Cap the request body so a huge upload can't exhaust memory/disk.
	r.Body = http.MaxBytesReader(w, r.Body, upload.MaxImageBytes)
	if err := r.ParseMultipartForm(upload.MaxImageBytes); err != nil {
		web.RespondError(w, resource, apperror.Validation("file is too large (max 5MB) or the form is invalid"))
		return "", false
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		web.RespondError(w, resource, apperror.Validation("missing 'file' field"))
		return "", false
	}
	defer file.Close()

	relPath, err := upload.SaveImage(uploadDir, "avatars", file)
	if err != nil {
		if errors.Is(err, upload.ErrUnsupportedType) {
			web.RespondError(w, resource, apperror.Validation("file must be a JPEG, PNG or WebP image"))
		} else {
			web.RespondError(w, resource, err)
		}
		return "", false
	}
	return relPath, true
}
