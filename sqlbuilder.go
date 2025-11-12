// the stmt placeholder using '?' for all, it will be replaced by builder.
package qsql

import (
	"fmt"
	"reflect"
	"strings"
)

type SqlBuilder struct {
	drvName string

	queryStr string
	fromBuff strings.Builder

	args []interface{}

	indent string

	dump bool
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

func (b *SqlBuilder) SetDump(dump bool) {
	b.dump = dump
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
		dump:     b.dump,
	}
	copy(n.args, b.args)
	n.fromBuff.WriteString(b.fromBuff.String())
	return n
}

// add the query to the buffer with indent,
// and add the args to the argument recorder, call builder.Args() to output the recorder
func (b *SqlBuilder) Add(query string, args ...interface{}) *SqlBuilder {
	if len(query) == 0 {
		return b
	}

	b.fromBuff.WriteString(query)
	b.fromBuff.WriteString(b.Indent())
	if len(args) > 0 {
		b.args = append(b.args, args...)
	}
	return b
}

// if indent isn't " ", add one tab with two space width to the buffer before adding
func (b *SqlBuilder) AddTab(query string, args ...interface{}) *SqlBuilder {
	if b.Indent() != " " {
		return b.Add("  "+query, args...)
	}
	return b.Add(query, args...)
}

// call AddTab if ok is true
func (b *SqlBuilder) AddIf(ok bool, query string, args ...interface{}) *SqlBuilder {
	if !ok {
		return b
	}
	return b.AddTab(query, args...)
}

// append the slice to the sql params and return then the stmt string.
// where in is not a slice kind, it will be panic
func (b *SqlBuilder) AddStmtIn(inArgs interface{}) string {
	v := reflect.ValueOf(inArgs)
	if v.Kind() != reflect.Slice {
		panic("StmtIn input is not a slice type")
	}
	if v.Len() == 0 {
		panic("need arguments of in condition")
	}
	stmtIn := make([]rune, v.Len()*2)
	args := make([]interface{}, v.Len())
	for i := v.Len() - 1; i > -1; i-- {
		stmtIn[i*2] = '?'
		stmtIn[i*2+1] = ','
		args[i] = v.Index(i).Interface()
	}
	b.args = append(b.args, args...)
	return string(stmtIn[:len(stmtIn)-1])
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
	sqlStr := ""
	if len(b.queryStr) > 0 {
		sqlStr = strings.TrimSuffix(b.queryStr+b.Indent()+b.fromBuff.String(), b.Indent())
	} else {
		sqlStr = strings.TrimSuffix(b.fromBuff.String(), b.Indent())
	}

	// fix build driver
	switch b.drvName {
	case DRV_NAME_ORACLE, _DRV_NAME_OCI8:
		paramIdx := 1
		buff := strings.Builder{}
		for _, r := range sqlStr {
			if r != '?' {
				buff.WriteRune(r)
			} else {
				buff.WriteString(fmt.Sprintf(":%d", paramIdx))
				paramIdx++
			}
		}
		sqlStr = buff.String()
	case DRV_NAME_POSTGRES:
		paramIdx := 1
		buff := strings.Builder{}
		for _, r := range sqlStr {
			if r != '?' {
				buff.WriteRune(r)
			} else {
				buff.WriteString(fmt.Sprintf("$%d", paramIdx))
				paramIdx++
			}
		}
		sqlStr = buff.String()
	case DRV_NAME_SQLSERVER, _DRV_NAME_MSSQL:
		buff := strings.Builder{}
		paramIdx := 1
		for _, r := range sqlStr {
			if r != '?' {
				buff.WriteRune(r)
			} else {
				buff.WriteString(fmt.Sprintf("@p%d", paramIdx))
				paramIdx++
			}
		}
		sqlStr = buff.String()
	default:
		// nothing to do.
	}
	if b.dump {
		log.Println(sqlStr, b.args)
	}
	return sqlStr
}

func (b *SqlBuilder) Args() []interface{} {
	return b.args
}

func (b *SqlBuilder) Sql() []interface{} {
	result := []interface{}{b.String()}
	return append(result, b.args...)
}
