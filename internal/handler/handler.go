// Package handler contains the HTTP layer: it decodes requests, delegates
// to the service layer and writes responses. It holds no business logic.
package handler

import (
	"net/http"
	"strconv"

	"story-go-mysql/internal/apperror"
	"story-go-mysql/internal/model"
	"story-go-mysql/internal/web"
)

// parseID reads the {id} path value and parses it as an unsigned integer.
// On failure it writes a 400 response and reports ok=false so the caller
// can return early.
func parseID(w http.ResponseWriter, r *http.Request, resource string) (uint64, bool) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		web.RespondError(w, resource, apperror.Validation("invalid "+resource+" id"))
		return 0, false
	}
	return id, true
}

// parseListParams reads q/page/pageSize from the query string. Missing or
// non-numeric page/pageSize stay at 0, which the service normalizes to
// sensible defaults.
func parseListParams(r *http.Request) model.ListParams {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	return model.ListParams{Query: q.Get("q"), Page: page, PageSize: pageSize}
}
