package qsql

import (
	"fmt"
	"testing"
)

func TestSqlBuilder(t *testing.T) {
	sb := NewSqlBuilder(DRV_NAME_POSTGRES)
	sb.Add("SELECT")
	sb.AddTab(sb.Select("tb1.id", "count(*)"))
	sb.Add("FROM")
	sb.AddTab("tmp tb1")
	sb.AddTab("INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id")
	sb.Add("WHERE")
	sb.AddTab("1=1")
	sb.AddTabOk(false, "AND (1=?)", 0)
	sb.AddTab("OR (tb1 IN (" + sb.In([]interface{}{1, 2}) + "))")
	sb.Add("GROUP BY tb1.id")
	sb.Add("HAVING count(*)>?", 1)
	sb.Add("ORDER BY tb1.id DESC")
	sb.Add("OFFSET ?", 1)
	sb.Add("LIMIT ?", 1)
	fmt.Println(sb)

	sb1 := NewSqlBuilder(DRV_NAME_POSTGRES)
	sb1.Add("SELECT * FROM ("+sb.String()+") tmp", sb1.Args())
	fmt.Println(sb1)
}
