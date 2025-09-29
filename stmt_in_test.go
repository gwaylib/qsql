package qsql

import "testing"

func TestStmtIn(t *testing.T) {
	sqlite3Output := StmtIn(1, 3, DRV_NAME_SQLITE3)
	if sqlite3Output != "?,?,?" {
		t.Fatalf("expect '?,?,?', but: %s", sqlite3Output)
	}
	pgOutput := StmtIn(1, 3, DRV_NAME_POSTGRES)
	if pgOutput != "$1,$2,$3" {
		t.Fatalf("expect '$1,$2,$3', but: %s", pgOutput)
	}
	pgOutput1 := StmtIn(2, 3, DRV_NAME_POSTGRES)
	if pgOutput1 != "$2,$3,$4" {
		t.Fatalf("expect '$2,$3,$4', but: %s", pgOutput1)
	}
	msOutput := StmtIn(0, 3, DRV_NAME_SQLSERVER)
	if msOutput != "@p0,@p1,@p2" {
		t.Fatalf("expect '@p0,@p1,@p2', but@ %s", msOutput)
	}
	msOutput1 := StmtIn(1, 3, DRV_NAME_SQLSERVER)
	if msOutput1 != "@p1,@p2,@p3" {
		t.Fatalf("expect '@p1,@p2,@p3', but@ %s", msOutput1)
	}
}
