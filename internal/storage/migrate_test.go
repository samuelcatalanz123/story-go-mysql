package storage

import "testing"

func TestSplitStatements(t *testing.T) {
	in := "CREATE TABLE a (id INT);\n\nCREATE TABLE b (id INT);\n"
	got := splitStatements(in)
	if len(got) != 2 {
		t.Fatalf("esperaba 2 sentencias, obtuve %d: %v", len(got), got)
	}
	if got[0] != "CREATE TABLE a (id INT)" {
		t.Fatalf("sentencia 0 inesperada: %q", got[0])
	}
	if got[1] != "CREATE TABLE b (id INT)" {
		t.Fatalf("sentencia 1 inesperada: %q", got[1])
	}
}

func TestSplitStatementsIgnoresEmpty(t *testing.T) {
	if got := splitStatements("   ;\n;  "); len(got) != 0 {
		t.Fatalf("esperaba 0 sentencias, obtuve %v", got)
	}
}
