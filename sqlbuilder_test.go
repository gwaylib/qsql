package qsql

import (
	"fmt"
	"testing"
)

func TestSqlBuilder(t *testing.T) {
	sb := NewSqlBuilder("tmp tb1")
	fmt.Println(sb.Select("*"))

	sb.Joins("INNER JOIN tmp1 tb2 ON tb2.id=tb2.tmp_id")
	fmt.Println(sb.Select("*"))

	sb.Where("1=?", 0)
	fmt.Println(sb.Select("*"))

	sb.WhereOr("1=1")
	fmt.Println(sb.Select("*"))

	sb.GroupBy("tb1.id")
	fmt.Println(sb.Select("*"))

	sb.GroupBy("tb2.id")
	fmt.Println(sb.Select("*"))

	sb.Having("count(tb1.id)>?", 2)
	fmt.Println(sb.Select("*"))

	sb.Offset(0)
	fmt.Println(sb.Select("*"))
	sb.Offset(1)
	fmt.Println(sb.Select("*"))

	sb.Limit(0)
	fmt.Println(sb.Select("*"))
	sb.Limit(1)
	fmt.Println(sb.Select("*"))

	mTable, args := sb.Select("id")
	fmt.Println(NewSqlBuilder("("+mTable+") tmp", args...).Select("*"))
}
