package qsql

import (
	"strings"
)

type SqlBuilder struct {
	drvName string

	queryStr string
	fromBuff strings.Builder

	args []interface{}

	indent string
}

func NewSqlBuilder(drvName ...string) *SqlBuilder {
	b := &SqlBuilder{
		indent: " ",
	}
	if len(drvName) > 0 {
		b.drvName = drvName[0]
	} else {
		b.drvName = DRV_NAME_SQLITE3
	}
	return b
}

func (b *SqlBuilder) Copy() *SqlBuilder {
	n := &SqlBuilder{
		drvName:  b.drvName,
		queryStr: b.queryStr,
		args:     make([]interface{}, len(b.args)),
		indent:   b.indent,
	}
	copy(n.args, b.args)
	n.fromBuff.WriteString(b.fromBuff.String())
	return n
}
func (b *SqlBuilder) SetIndent(indent string) *SqlBuilder {
	b.indent = indent
	return b
}

func (b *SqlBuilder) Add(key string, args ...interface{}) *SqlBuilder {
	if len(key) == 0 {
		return b
	}

	b.fromBuff.WriteString(key)
	b.fromBuff.WriteString(b.indent)
	if len(args) > 0 {
		b.args = append(b.args, args...)
	}
	return b
}

// recursive Add,  only format the code when coding
func (b *SqlBuilder) AddTab(key string, args ...interface{}) *SqlBuilder {
	return b.Add(key, args...)
}

func (b *SqlBuilder) AddIf(ok bool, key string, args ...interface{}) *SqlBuilder {
	if !ok {
		return b
	}
	return b.Add(key, args...)
}

func (b *SqlBuilder) AddIn(in []interface{}) string {
	if len(in) == 0 {
		panic("need arguments of in condition")
	}
	b.args = append(b.args, in...)
	return stmtIn(len(b.args)-1, len(in), b.drvName)
}

func (b *SqlBuilder) Select(column ...string) *SqlBuilder {
	b.queryStr = "SELECT" + b.indent
	if len(column) > 0 {
		b.queryStr += strings.Join(column, ",")
	} else {
		b.queryStr += "*"
	}
	return b
}
func (b *SqlBuilder) SelectStruct(obj interface{}) *SqlBuilder {
	fields, err := reflectInsertStruct(obj, b.drvName)
	if err != nil {
		panic(err)
	}
	return b.Select(fields.Names...)
}

func (b *SqlBuilder) String() string {
	if len(b.queryStr) > 0 {
		return strings.TrimSuffix(b.queryStr+b.indent+b.fromBuff.String(), b.indent)
	}
	return strings.TrimSuffix(b.fromBuff.String(), b.indent)
}

func (b *SqlBuilder) Args() []interface{} {
	return b.args
}

func (b *SqlBuilder) Sql() []interface{} {
	result := []interface{}{b.String()}
	return append(result, b.args...)
}
