package qsql

import (
	"testing"
)

func TestSelectBuilder(t *testing.T) {
	paramIn := []int{1, 2}
	bd := NewSelectBuilder(DRV_NAME_POSTGRES)
	bd.Select("count(*)")
	bd.From("tmp tb1")
	bd.From("INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id")
	bd.Where("1=1")
	bd.Where("AND (1=?)", 0)
	bd.WhereIn("OR (tb1 IN ?)", paramIn)
	bd.GroupBy("tb1.id")
	bd.GroupBy("HAVING count(*)>?", 1)
	if bd.String() !=
		`SELECT count(*) FROM tmp tb1 INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id WHERE 1=1 AND (1=$1) OR (tb1 IN ($2,$3)) GROUP BY tb1.id HAVING count(*)>$4` {
		t.Fatal(bd)
	}

	bd1 := bd.Copy(true)
	bd1.OrderBy("tb1.id DESC")
	bd1.Offset(1)
	bd1.Limit(1)
	bd1.Select("tb1.id", "count(*)")
	if bd1.String() !=
		`SELECT tb1.id, count(*) FROM tmp tb1 INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id WHERE 1=1 AND (1=$1) OR (tb1 IN ($2,$3)) GROUP BY tb1.id HAVING count(*)>$4 ORDER BY tb1.id DESC LIMIT 1 OFFSET 1` {
		t.Fatal(bd1)
	}

	bd2 := NewSelectBuilder(DRV_NAME_POSTGRES)
	bd2.Select("*")
	bd2.From("("+bd.String()+") tmp", bd1.Args()...)
	if bd2.String() !=
		`SELECT * FROM (SELECT count(*) FROM tmp tb1 INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id WHERE 1=1 AND (1=$1) OR (tb1 IN ($2,$3)) GROUP BY tb1.id HAVING count(*)>$4) tmp` {
		t.Fatal(bd2)
	}

	if len(bd2.Sql()) != 5 {
		t.Fatalf("%+v", bd2.Sql())
	}
}
