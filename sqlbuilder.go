package qsql

import (
	"strings"
)

type SqlBuilder struct {
	drvName string
	buff    strings.Builder
	args    []interface{}
}

func NewSqlBuilder(drvName ...string) *SqlBuilder {
	b := &SqlBuilder{}
	if len(drvName) > 0 {
		b.drvName = drvName[0]
	} else {
		b.drvName = DRV_NAME_SQLITE3
	}
	return b
}

func (b *SqlBuilder) Add(key string, args ...interface{}) *SqlBuilder {
	if len(key) > 0 {
		b.buff.WriteString(key)
	}
	if len(args) > 0 {
		b.args = append(b.args, args...)
	}
	b.buff.WriteString("\n")
	return b
}

func (b *SqlBuilder) AddTab(key string, args ...interface{}) *SqlBuilder {
	b.buff.WriteString("  ")
	return b.Add(key, args)
}
func (b *SqlBuilder) AddTabOK(ok bool, key string, args ...interface{}) *SqlBuilder {
	if !ok {
		return b
	}
	return b.AddTab(key, args...)
}

func (b *SqlBuilder) In(in []interface{}) string {
	if len(in) == 0 {
		panic("need condition")
	}
	b.args = append(b.args, in...)
	return stmtIn(len(b.args)-1, len(in), b.drvName)
}

func (b *SqlBuilder) Args() []interface{} {
	return b.args
}

func (b *SqlBuilder) Select(column ...string) string {
	selectStr := "SELECT\n  "
	if len(column) > 0 {
		selectStr += strings.Join(column, ", ")
	} else {
		selectStr += "*"
	}
	return selectStr + "\n" + b.buff.String()
}
func (b *SqlBuilder) SelectStruct(obj interface{}) string {
	fields, err := reflectInsertStruct(obj, b.drvName)
	if err != nil {
		panic(err)
	}
	return "SELECT\n  " + fields.Names + "\n" + b.buff.String()
}
