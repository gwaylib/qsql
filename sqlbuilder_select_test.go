package qsql

import (
	"testing"
)

func TestSelectBuilder(t *testing.T) {
	paramIn := []int{1, 2}
	bd := NewSelectBuilder()
	bd.Select("count(*)")
	bd.From("tmp tb1")
	bd.From("INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id")
	bd.Where(true, "1=1")
	bd.Where(true, "AND (1=?)", 0)
	bd.Where(true, "OR (tb1 IN ("+bd.In(paramIn)+"))")
	bd.Group("tb1.id")
	bd.Group("HAVING count(*)>?", 1)
	if bd.StrTo(DRV_NAME_POSTGRES) !=
		`SELECT count(*) FROM tmp tb1 INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id WHERE 1=1 AND (1=$1) OR (tb1 IN ($2,$3)) GROUP BY tb1.id HAVING count(*)>$4` {
		t.Fatal(bd)
	}

	bd1 := bd.Copy()
	bd1.Order("tb1.id DESC")
	bd1.Offset(1)
	bd1.Limit(1)
	bd1.Select("tb1.id", "count(*)")
	if bd1.StrTo(DRV_NAME_POSTGRES) !=
		`SELECT tb1.id, count(*) FROM tmp tb1 INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id WHERE 1=1 AND (1=$1) OR (tb1 IN ($2,$3)) GROUP BY tb1.id HAVING count(*)>$4 ORDER BY tb1.id DESC OFFSET 1 LIMIT 1` {
		t.Fatal(bd1)
	}

	bd2 := NewSelectBuilder()
	bd2.Select("*")
	bd2.From("("+bd.StrTo(DRV_NAME_POSTGRES)+") tmp", bd1.Args()...)
	if bd2.StrTo(DRV_NAME_POSTGRES) !=
		`SELECT * FROM (SELECT count(*) FROM tmp tb1 INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id WHERE 1=1 AND (1=$1) OR (tb1 IN ($2,$3)) GROUP BY tb1.id HAVING count(*)>$4) tmp` {
		t.Fatal(bd2)
	}

	if len(bd2.SqlTo(DRV_NAME_POSTGRES)) != 5 {
		// [SELECT * FROM (SELECT count(*) FROM tmp tb1 INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id WHERE 1=1 AND (1=$1) OR (tb1 IN ($2,$3)) GROUP BY tb1.id HAVING count(*)>$4) tmp 0 1 2 1]
		t.Fatalf("%+v", bd2.SqlTo(DRV_NAME_POSTGRES))
	}
}

// func TestSqlBuilderUpdate(t *testing.T) {
// 	bd := NewSqlBuilder(DRV_NAME_POSTGRES)
// 	bd.Update("UPDATE tmp SET")
// 	bd.TabAdd("(val1=?, val2=?)", 1, 2)
// 	bd.Add("WHERE")
// 	bd.TabAdd("id=?", 1)
// 	if bd.String() !=
// 		`UPDATE tmp SET (val1=$1, val2=$2) WHERE id=$3` {
// 		t.Fatal(bd)
// 	}
// 	if len(bd.Sql()) != 4 {
// 		t.Fatal(bd.Sql())
// 	}
//
// 	bd1 := NewSqlBuilder(DRV_NAME_POSTGRES)
// 	bd1.Add("UPDATE tmp SET")
// 	bd1.TabAdd("val1=?,", 1)
// 	bd1.TabAdd("val2=?", 2)
// 	bd1.Add("WHERE")
// 	bd1.TabAdd("id=?", 1)
// 	if bd1.String() !=
// 		`UPDATE tmp SET val1=$1, val2=$2 WHERE id=$3` {
// 		t.Fatal(bd1)
// 	}
// 	if len(bd1.Sql()) != 4 {
// 		t.Fatal(bd1.Sql())
// 	}
// }
//
// func TestSqlBuilderDel(t *testing.T) {
// 	bd := NewSqlBuilder(DRV_NAME_POSTGRES)
// 	bd.Add("DELETE FROM tmp WHERE")
// 	bd.TabAdd("id=?", 1)
// 	if bd.String() !=
// 		`DELETE FROM tmp WHERE id=$1` {
// 		t.Fatal(bd)
// 	}
// 	if len(bd.Sql()) != 2 {
// 		t.Fatal(bd.Sql())
// 	}
// }
