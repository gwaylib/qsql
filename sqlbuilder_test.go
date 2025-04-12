package qsql

import (
	"testing"
)

func TestSqlBuilder(t *testing.T) {
	bd := NewSqlBuilder(DRV_NAME_POSTGRES)
	bd.Select("count(*)")
	bd.Add("FROM")
	bd.AddTab("tmp tb1")
	bd.AddTab("INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id")
	bd.Add("WHERE")
	bd.AddTab("1=1")
	bd.AddIf(true, "AND (1=?)", 0)
	bd.AddIf(true, "OR (tb1 IN ("+bd.AddIn([]interface{}{1, 2})+"))")
	bd.Add("GROUP BY tb1.id")
	bd.Add("HAVING count(*)>?", 1)
	if bd.String() !=
		`SELECT count(*) FROM tmp tb1 INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id WHERE 1=1 AND (1=?) OR (tb1 IN (:2,:3)) GROUP BY tb1.id HAVING count(*)>?` {
		t.Fatal(bd)
	}

	bd1 := bd.Copy()
	bd1.Add("ORDER BY tb1.id DESC")
	bd1.Add("OFFSET ?", 1)
	bd1.Add("LIMIT ?", 1)
	bd1.Select("tb1.id", "count(*)")
	if bd1.String() !=
		`SELECT tb1.id, count(*) FROM tmp tb1 INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id WHERE 1=1 AND (1=?) OR (tb1 IN (:2,:3)) GROUP BY tb1.id HAVING count(*)>? ORDER BY tb1.id DESC OFFSET ? LIMIT ?` {
		t.Fatal(bd1)
	}

	bd2 := NewSqlBuilder(DRV_NAME_POSTGRES)
	bd2.Select("*")
	bd2.Add("FROM ("+bd.String()+") tmp", bd1.Args()...)
	if bd2.String() !=
		`SELECT * FROM (SELECT count(*) FROM tmp tb1 INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id WHERE 1=1 AND (1=?) OR (tb1 IN (:2,:3)) GROUP BY tb1.id HAVING count(*)>?) tmp` {
		t.Fatal(bd2)
	}

	if len(bd2.Sql()) != 7 {
		t.Fatal(bd2.Sql())
	}
}
