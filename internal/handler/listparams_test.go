package handler

import (
	"net/http/httptest"
	"testing"
)

func TestParseListParamsDefaults(t *testing.T) {
	p := parseListParams(httptest.NewRequest("GET", "/characters", nil))
	if p.Query != "" || p.Page != 0 || p.PageSize != 0 {
		t.Fatalf("inesperado: %+v", p)
	}
}

func TestParseListParamsReadsValues(t *testing.T) {
	p := parseListParams(httptest.NewRequest("GET", "/characters?q=asha&page=2&pageSize=5", nil))
	if p.Query != "asha" || p.Page != 2 || p.PageSize != 5 {
		t.Fatalf("inesperado: %+v", p)
	}
}
