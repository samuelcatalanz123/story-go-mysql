package repository

import "testing"

func TestBuildSearchEmpty(t *testing.T) {
	clause, args := buildSearch("")
	if clause != "" {
		t.Fatalf("clause=%q, esperaba vacío", clause)
	}
	if args != nil {
		t.Fatalf("args=%v, esperaba nil", args)
	}
}

func TestBuildSearchNonEmpty(t *testing.T) {
	clause, args := buildSearch("asha")
	if clause == "" {
		t.Fatal("esperaba una cláusula WHERE")
	}
	if len(args) != 2 || args[0] != "%asha%" || args[1] != "%asha%" {
		t.Fatalf("args=%v, esperaba dos patrones %%asha%%", args)
	}
}
