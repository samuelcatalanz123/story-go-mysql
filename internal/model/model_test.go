package model

import "testing"

func TestListParamsNormalize(t *testing.T) {
	got := ListParams{Page: 0, PageSize: 999}.Normalize()
	if got.Page != 1 {
		t.Fatalf("page=%d, esperaba 1", got.Page)
	}
	if got.PageSize != 100 {
		t.Fatalf("pageSize=%d, esperaba 100", got.PageSize)
	}
	if d := (ListParams{}).Normalize(); d.PageSize != 20 {
		t.Fatalf("default pageSize=%d, esperaba 20", d.PageSize)
	}
}

func TestListParamsLimitOffset(t *testing.T) {
	p := ListParams{Page: 3, PageSize: 20}.Normalize()
	if p.Limit() != 20 || p.Offset() != 40 {
		t.Fatalf("limit=%d offset=%d, esperaba 20 y 40", p.Limit(), p.Offset())
	}
}
