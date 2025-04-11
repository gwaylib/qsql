package qsql

import (
	"fmt"
	"testing"
)

func TestSqlBuilder(t *testing.T) {
	sb := NewSqlBuilder(DRV_NAME_POSTGRES)
	sb.Add("FROM")
	sb.AddTab("tmp tb1")
	sb.AddTab("INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id")
	sb.Add("WHERE")
	sb.AddTab("1=1")
	sb.AddTabOK(false, "AND (1=?)", 0)
	sb.AddTabOK(true, "OR (tb1 IN ("+sb.In([]interface{}{1, 2})+"))")
	sb.Add("GROUP BY tb1.id")
	sb.Add("HAVING count(*)>?", 1)
	fmt.Println(sb.Select("count(*)"))

	sb.Add("ORDER BY tb1.id DESC")
	sb.Add("OFFSET ?", 1)
	sb.Add("LIMIT ?", 1)
	fmt.Println(sb.Select("tb1.id", "count(*)"))

	sb1 := NewSqlBuilder(DRV_NAME_POSTGRES)
	sb1.Add("FROM ("+sb.Select("tb1.id", "count(*)")+") tmp", sb.Args())
	fmt.Println(sb1.Select("*"))
}
