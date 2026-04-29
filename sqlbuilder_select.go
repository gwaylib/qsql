// the stmt placeholder using '?' for all, it will be replaced by builder.
package qsql

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type SelectBuilder struct {
	indent string
	dump   bool

	queryStr  string
	fromStr   string
	fromArgs  []interface{}
	whereStr  string
	whereArgs []interface{}
	groupStr  string
	groupArgs []interface{}
	orderStr  string
	orderArgs []interface{}
	offset    int64
	limit     int64
}

func NewSelectBuilder() *SelectBuilder {
	return NewSelectBuilderWithIndent(" ")
}

func NewSelectBuilderWithIndent(indent string, drvName ...string) *SelectBuilder {
	b := &SelectBuilder{
		indent: indent,
	}
	return b
}

func (b *SelectBuilder) SetDump(dump bool) *SelectBuilder {
	b.dump = dump
	return b
}

func (b *SelectBuilder) Indent() string {
	if len(b.indent) == 0 {
		return " "
	}
	return b.indent
}

func (b *SelectBuilder) Copy() *SelectBuilder {
	n := &SelectBuilder{
		indent:    b.indent,
		dump:      b.dump,
		queryStr:  b.queryStr,
		fromStr:   b.fromStr,
		fromArgs:  make([]interface{}, len(b.fromArgs)),
		whereStr:  b.whereStr,
		whereArgs: make([]interface{}, len(b.whereArgs)),
		groupStr:  b.groupStr,
		groupArgs: make([]interface{}, len(b.groupArgs)),
		orderStr:  b.orderStr,
		orderArgs: make([]interface{}, len(b.orderArgs)),
	}
	copy(n.fromArgs, b.fromArgs)
	copy(n.whereArgs, b.whereArgs)
	copy(n.groupArgs, b.groupArgs)
	copy(n.orderArgs, b.orderArgs)
	return n
}

func (b *SelectBuilder) From(query string, args ...interface{}) *SelectBuilder {
	if len(query) == 0 {
		return b
	}
	if len(b.fromStr) > 0 {
		b.fromStr += b.Indent()
	}
	b.fromStr += query
	if len(args) > 0 {
		b.fromArgs = append(b.fromArgs, args...)
	}
	return b
}

func (b *SelectBuilder) Where(ok bool, query string, args ...interface{}) *SelectBuilder {
	if !ok {
		return b
	}
	if len(query) == 0 {
		return b
	}
	if len(b.whereStr) > 0 {
		b.whereStr += b.Indent()
	}
	b.whereStr += query
	if len(args) > 0 {
		b.whereArgs = append(b.whereArgs, args...)
	}
	return b
}

// append the slice to the sql params and return then the stmt string.
// where in is not a slice kind, it will be panic
func (b *SelectBuilder) In(inArgs interface{}) string {
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
	b.whereArgs = append(b.whereArgs, args...)
	return string(stmtIn[:len(stmtIn)-1])
}

func (b *SelectBuilder) Group(query string, args ...interface{}) *SelectBuilder {
	if len(query) == 0 {
		return b
	}
	if len(b.groupStr) > 0 {
		b.groupStr += b.Indent()
	}
	b.groupStr += query
	if len(args) > 0 {
		b.groupArgs = append(b.groupArgs, args...)
	}
	return b
}

func (b *SelectBuilder) Order(query string, args ...interface{}) *SelectBuilder {
	if len(query) == 0 {
		return b
	}
	if len(b.orderStr) > 0 {
		b.orderStr += ("," + b.Indent())
	}
	b.orderStr += query
	if len(args) > 0 {
		b.orderArgs = append(b.orderArgs, args...)
	}
	return b
}

func (b *SelectBuilder) Offset(offset int64) *SelectBuilder {
	b.offset = offset
	return b
}
func (b *SelectBuilder) Limit(limit int64) *SelectBuilder {
	b.limit = limit
	return b
}

func (b *SelectBuilder) Select(column ...string) *SelectBuilder {
	if len(column) > 0 {
		b.queryStr = strings.Join(column, ", ")
	} else {
		b.queryStr = "*"
	}
	return b
}
func (b *SelectBuilder) SelectStruct(obj interface{}) *SelectBuilder {
	fields, err := reflectSelectStruct(obj)
	if err != nil {
		panic(err)
	}
	return b.Select(fields...)
}

func (b *SelectBuilder) Args() []interface{} {
	result := make([]interface{}, len(b.fromArgs)+len(b.whereArgs)+len(b.groupArgs)+len(b.orderArgs))
	idx := copy(result, b.fromArgs)
	idx += copy(result[idx:], b.whereArgs)
	idx += copy(result[idx:], b.groupArgs)
	idx += copy(result[idx:], b.orderArgs)
	return result
}

func (b *SelectBuilder) StrTo(drvName string) string {
	if len(b.queryStr) == 0 {
		b.queryStr = "*"
	}
	sqlStr := "SELECT " + b.queryStr
	if len(b.fromStr) > 0 {
		sqlStr += (b.Indent() + "FROM " + b.fromStr)
	}
	if len(b.whereStr) > 0 {
		sqlStr += (b.Indent() + "WHERE " + b.whereStr)
	}
	if len(b.groupStr) > 0 {
		sqlStr += (b.Indent() + "GROUP BY " + b.groupStr)
	}
	if len(b.orderStr) > 0 {
		sqlStr += (b.Indent() + "ORDER BY " + b.orderStr)
	}
	if b.offset > 0 {
		sqlStr += (b.Indent() + "OFFSET " + strconv.FormatInt(b.offset, 10))
	}
	if b.limit > 0 {
		sqlStr += (b.Indent() + "LIMIT " + strconv.FormatInt(b.limit, 10))
	}

	// fix build driver
	switch drvName {
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
		log.Println(sqlStr, b.Args())
	}
	return sqlStr
}

func (b *SelectBuilder) SqlTo(drvName string) []interface{} {
	result := []interface{}{b.StrTo(drvName)}
	return append(result, b.Args()...)
}
