package qsql

import (
	"reflect"
	"strings"
)

type SqlBuilder struct {
	drvName string

	queryStr string
	fromBuff strings.Builder

	args []interface{}

	indent string
}

func NewSqlBuilderWithIndent(indent string, drvName ...string) *SqlBuilder {
	b := &SqlBuilder{
		indent: indent,
	}
	if len(drvName) > 0 {
		b.drvName = drvName[0]
	}
	return b
}

func NewSqlBuilder(drvName ...string) *SqlBuilder {
	return NewSqlBuilderWithIndent(" ", drvName...)
}

func (b *SqlBuilder) DrvName() string {
	if len(b.drvName) > 0 {
		return b.drvName
	}
	return DRV_NAME_SQLITE3
}

func (b *SqlBuilder) Indent() string {
	if len(b.indent) == 0 {
		return " "
	}
	return b.indent
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

// add the key to the buffer with indent,
// and append the args to the argument recorder, call builder.Args() to output the recorder
func (b *SqlBuilder) Add(key string, args ...interface{}) *SqlBuilder {
	if len(key) == 0 {
		return b
	}

	b.fromBuff.WriteString(key)
	b.fromBuff.WriteString(b.Indent())
	if len(args) > 0 {
		b.args = append(b.args, args...)
	}
	return b
}

// if indent isn't " ", add one tab with two space width to the buffer before adding
func (b *SqlBuilder) AddTab(key string, args ...interface{}) *SqlBuilder {
	if b.Indent() != " " {
		return b.Add("  "+key, args...)
	}
	return b.Add(key, args...)
}

// call AddTab if ok is true
func (b *SqlBuilder) AddIf(ok bool, key string, args ...interface{}) *SqlBuilder {
	if !ok {
		return b
	}
	return b.AddTab(key, args...)
}

// where in is not a slice kind, it will be panic
func (b *SqlBuilder) AddStmtIn(in interface{}) string {
	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Slice {
		panic("StmtIn input is not a slice type")
	}
	if v.Len() == 0 {
		panic("need arguments of in condition")
	}
	args := make([]interface{}, v.Len())
	for i := v.Len() - 1; i > -1; i-- {
		args[i] = v.Index(i).Interface()
	}
	stmtIn := stmtIn(len(b.args), len(args), b.drvName)
	b.args = append(b.args, args...)
	return stmtIn
}

func (b *SqlBuilder) Select(column ...string) *SqlBuilder {
	b.queryStr = "SELECT" + b.Indent()
	if len(column) > 0 {
		b.queryStr += strings.Join(column, ", ")
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
		return strings.TrimSuffix(b.queryStr+b.Indent()+b.fromBuff.String(), b.Indent())
	}
	return strings.TrimSuffix(b.fromBuff.String(), b.Indent())
}

func (b *SqlBuilder) Args() []interface{} {
	return b.args
}

func (b *SqlBuilder) Sql() []interface{} {
	result := []interface{}{b.String()}
	return append(result, b.args...)
}
