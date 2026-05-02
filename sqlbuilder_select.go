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

	queryStr    string
	fromStr     string
	fromArgs    []interface{}
	whereStr    string
	whereArgs   []interface{}
	groupByStr  string
	groupByArgs []interface{}
	orderByStr  string
	orderByArgs []interface{}
	offset      int64
	limit       int64
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

// copy and return a new selector without select buffer
func (b *SelectBuilder) Copy(newSelectBuffer bool) *SelectBuilder {
	queryStr := ""
	if !newSelectBuffer {
		queryStr = b.queryStr
	}
	n := &SelectBuilder{
		indent:      b.indent,
		dump:        b.dump,
		queryStr:    queryStr,
		fromStr:     b.fromStr,
		fromArgs:    make([]interface{}, len(b.fromArgs)),
		whereStr:    b.whereStr,
		whereArgs:   make([]interface{}, len(b.whereArgs)),
		groupByStr:  b.groupByStr,
		groupByArgs: make([]interface{}, len(b.groupByArgs)),
		orderByStr:  b.orderByStr,
		orderByArgs: make([]interface{}, len(b.orderByArgs)),
		offset:      b.offset,
		limit:       b.limit,
	}
	copy(n.fromArgs, b.fromArgs)
	copy(n.whereArgs, b.whereArgs)
	copy(n.groupByArgs, b.groupByArgs)
	copy(n.orderByArgs, b.orderByArgs)
	return n
}

// select the columns and append to select buffer
func (b *SelectBuilder) Select(column ...string) *SelectBuilder {
	if len(column) > 0 {
		queryStr := strings.Join(column, ", ")
		if len(b.queryStr) > 0 {
			b.queryStr += (", " + queryStr)
		} else {
			b.queryStr = queryStr
		}
	} else if len(b.queryStr) == 0 {
		b.queryStr = "*"
	}
	return b
}

// clean select buffer and select the new columns
func (b *SelectBuilder) SelectNew(column ...string) *SelectBuilder {
	bd := b.Copy(true)
	return bd.Select(column...)
}

// clean select buffer and select the struct columns
func (b *SelectBuilder) SelectStruct(obj interface{}) *SelectBuilder {
	fields, err := reflectSelectStruct(obj)
	if err != nil {
		panic(err)
	}
	return b.SelectNew(fields...)
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

func (b *SelectBuilder) GroupBy(query string, args ...interface{}) *SelectBuilder {
	if len(query) == 0 {
		return b
	}
	if len(b.groupByStr) > 0 {
		b.groupByStr += b.Indent()
	}
	b.groupByStr += query
	if len(args) > 0 {
		b.groupByArgs = append(b.groupByArgs, args...)
	}
	return b
}

func (b *SelectBuilder) OrderBy(query string, args ...interface{}) *SelectBuilder {
	if len(query) == 0 {
		return b
	}
	if len(b.orderByStr) > 0 {
		b.orderByStr += ("," + b.Indent())
	}
	b.orderByStr += query
	if len(args) > 0 {
		b.orderByArgs = append(b.orderByArgs, args...)
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

// get the buffer args
func (b *SelectBuilder) Args() []interface{} {
	result := make([]interface{}, len(b.fromArgs)+len(b.whereArgs)+len(b.groupByArgs)+len(b.orderByArgs))
	idx := copy(result, b.fromArgs)
	idx += copy(result[idx:], b.whereArgs)
	idx += copy(result[idx:], b.groupByArgs)
	idx += copy(result[idx:], b.orderByArgs)
	return result
}

// translate sql to db driver
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
	if len(b.groupByStr) > 0 {
		sqlStr += (b.Indent() + "GROUP BY " + b.groupByStr)
	}
	if len(b.orderByStr) > 0 {
		sqlStr += (b.Indent() + "ORDER BY " + b.orderByStr)
	}
	if b.offset > 0 {
		sqlStr += (b.Indent() + "OFFSET " + strconv.FormatInt(b.offset, 10))
	}
	if b.limit > 0 {
		sqlStr += (b.Indent() + "LIMIT " + strconv.FormatInt(b.limit, 10))
	}

	// translate to db driver
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

// Merge StrTo and Args to a finally slice
func (b *SelectBuilder) SqlTo(drvName string) []interface{} {
	result := []interface{}{b.StrTo(drvName)}
	return append(result, b.Args()...)
}
