package qsql

import "testing"

func TestStmtWhereIn(t *testing.T) {
	sqlite3Output := StmtWhereIn(0, 3, DRV_NAME_SQLITE3)
	if sqlite3Output != "?,?,?" {
		t.Fatalf("expect '?,?,?', but: %s", sqlite3Output)
	}
	pgOutput := StmtWhereIn(0, 3, DRV_NAME_POSTGRES)
	if pgOutput != ":0,:1,:2" {
		t.Fatalf("expect ':0,:1,:2', but: %s", pgOutput)
	}
	pgOutput1 := StmtWhereIn(1, 3, DRV_NAME_POSTGRES)
	if pgOutput1 != ":1,:2,:3" {
		t.Fatalf("expect ':1,:2,:3', but: %s", pgOutput1)
	}
	msOutput := StmtWhereIn(0, 3, DRV_NAME_SQLSERVER)
	if msOutput != "@p0,@p1,@p2" {
		t.Fatalf("expect '@p0,@p1,@p2', but@ %s", msOutput)
	}
	msOutput1 := StmtWhereIn(1, 3, DRV_NAME_SQLSERVER)
	if msOutput1 != "@p1,@p2,@p3" {
		t.Fatalf("expect '@p1,@p2,@p3', but@ %s", msOutput1)
	}
}

func TestStmtSliceArgs(t *testing.T) {
	in := []string{"a", "b", "c"}
	out := StmtSliceArgs(in)
	if len(out) != 3 {
		t.Fatalf("expect 3, but:%d", len(out))
	}
	arg1, ok := out[0].(string)
	if !ok {
		t.Fatal("expect string, but not")
	}
	if arg1 != "a" {
		t.Fatalf("expect 'a' , but: %s ", arg1)
	}
}
