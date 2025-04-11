package qsql

import (
	"strconv"
	"strings"
)

type SqlBuilder struct {
	table   []string
	initArg []interface{}

	whereStr   []string
	whereArg   []interface{}
	whereInStr []string
	whereInArg [][]interface{}

	// TODO: how to treat or
	whereOrStr   []string
	whereOrArg   []interface{}
	whereOrInStr []string
	whereOrInArg [][]interface{}

	groupByStr  []string
	havingStr   []string
	havingArg   []interface{}
	havingOrStr []string
	havingOrArg []interface{}
	orderByStr  []string
	offset      int
	limit       int
}

func NewSqlBuilder(table string, args ...interface{}) *SqlBuilder {
	return &SqlBuilder{
		table:   []string{table},
		initArg: args,
	}
}

func (b *SqlBuilder) Copy() *SqlBuilder {
	n := &SqlBuilder{
		table:       make([]string, len(b.table)),
		whereStr:    make([]string, len(b.whereStr)),
		whereArg:    make([]interface{}, len(b.whereArg)),
		whereOrStr:  make([]string, len(b.whereOrStr)),
		whereOrArg:  make([]interface{}, len(b.whereOrArg)),
		groupByStr:  make([]string, len(b.groupByStr)),
		havingStr:   make([]string, len(b.havingStr)),
		havingArg:   make([]interface{}, len(b.havingArg)),
		havingOrStr: make([]string, len(b.havingOrStr)),
		havingOrArg: make([]interface{}, len(b.havingOrArg)),
		orderByStr:  make([]string, len(b.orderByStr)),
		offset:      b.offset,
		limit:       b.limit,
	}
	copy(n.table, b.table)
	copy(n.whereStr, b.whereStr)
	copy(n.whereArg, b.whereArg)
	copy(n.whereOrStr, b.whereOrStr)
	copy(n.whereOrArg, b.whereOrArg)
	copy(n.groupByStr, b.groupByStr)
	copy(n.havingStr, b.havingStr)
	copy(n.havingArg, b.havingArg)
	copy(n.havingOrStr, b.havingOrStr)
	copy(n.havingOrArg, b.havingOrArg)
	n.orderByStr = b.orderByStr
	n.offset = b.offset
	n.limit = b.limit
	return n
}

func (b *SqlBuilder) Joins(table string) *SqlBuilder {
	b.table = append(b.table, table)
	return b
}

func (b *SqlBuilder) Where(cond string, args ...interface{}) *SqlBuilder {
	b.whereStr = append(b.whereStr, cond)
	b.whereArg = append(b.whereArg, args...)
	return b
}
func (b *SqlBuilder) WhereIn(column string, args []interface{}) *SqlBuilder {
	if len(args) == 0 {
		return b
	}
	b.whereInStr = append(b.whereInStr, column)
	b.whereInArg = append(b.whereInArg, args)
	return b
}

func (b *SqlBuilder) WhereOr(cond string, args ...interface{}) *SqlBuilder {
	b.whereOrStr = append(b.whereOrStr, cond)
	b.whereOrArg = append(b.whereOrArg, args...)
	return b
}

func (b *SqlBuilder) WhereOrIn(cond string, args []interface{}) *SqlBuilder {
	if len(args) == 0 {
		return b
	}
	b.whereOrInStr = append(b.whereOrInStr, cond)
	b.whereOrInArg = append(b.whereOrInArg, args)
	return b
}

func (b *SqlBuilder) GroupBy(cond ...string) *SqlBuilder {
	b.groupByStr = append(b.groupByStr, cond...)
	return b
}
func (b *SqlBuilder) Having(cond string, args ...interface{}) *SqlBuilder {
	b.havingStr = append(b.havingStr, cond)
	b.havingArg = append(b.havingArg, args...)
	return b
}
func (b *SqlBuilder) HavingOr(cond string, args ...interface{}) *SqlBuilder {
	b.havingOrStr = append(b.havingOrStr, cond)
	b.havingOrArg = append(b.havingOrArg, args...)
	return b
}

func (b *SqlBuilder) OrderBy(cond string) *SqlBuilder {
	b.orderByStr = append(b.orderByStr, cond)
	return b
}
func (b *SqlBuilder) Offset(offset int) *SqlBuilder {
	b.offset = offset
	return b
}
func (b *SqlBuilder) Limit(limit int) *SqlBuilder {
	b.limit = limit
	return b
}
func (b *SqlBuilder) buildWhere() (string, []interface{}) {
	args := []interface{}{}
	if len(b.whereStr) == 0 && len(b.whereOrStr) == 0 {
		return "", args
	}

	query := "WHERE "
	if len(b.whereStr) > 0 {
		query += "(" + strings.Join(b.whereStr, " AND ") + ") "
		args = append(args, b.whereArg...)
	}
	if len(b.whereOrStr) > 0 {
		if len(b.whereStr) > 0 {
			query += "OR "
		}
		query += "(" + strings.Join(b.whereOrStr, " OR ") + ") "
		args = append(args, b.whereOrArg...)
	}
	return query, args
}
func (b *SqlBuilder) buildGroupBy() string {
	if len(b.groupByStr) == 0 {
		return ""
	}
	return "GROUP BY " + strings.Join(b.groupByStr, ", ") + " "
}
func (b *SqlBuilder) buildHaving() (string, []interface{}) {
	args := []interface{}{}
	if len(b.havingStr) == 0 && len(b.havingOrStr) == 0 {
		return "", args
	}
	query := "HAVING "
	if len(b.havingStr) > 0 {
		query += "(" + strings.Join(b.havingStr, " AND ") + ") "
		args = append(args, b.havingArg...)
	}
	if len(b.havingOrStr) > 0 {
		if len(b.havingStr) > 0 {
			query += "OR "
		}
		query += "(" + strings.Join(b.havingOrStr, " OR ") + ") "
		args = append(args, b.havingOrArg...)
	}
	return query, args
}
func (b *SqlBuilder) buildOrderBy() string {
	if len(b.orderByStr) == 0 {
		return ""
	}
	return "ORDER BY " + strings.Join(b.orderByStr, ", ") + " "
}
func (b *SqlBuilder) buildOffset() string {
	if b.offset == 0 {
		return ""
	}
	return "OFFSET " + strconv.Itoa(b.offset) + " "
}
func (b *SqlBuilder) buildLimit() string {
	if b.limit == 0 {
		return ""
	}
	return "LIMIT " + strconv.Itoa(b.limit) + " "
}

func (b *SqlBuilder) Select(column ...string) (string, []interface{}) {
	if len(b.table) == 0 {
		panic("table not set")
	}
	selectStr := "*"
	if len(column) > 0 {
		selectStr = strings.Join(column, ", ")
	}
	args := b.initArg
	query := "SELECT " + selectStr + " "
	query += "FROM " + strings.Join(b.table, " ") + " "

	whereStr, whereArg := b.buildWhere()
	if len(whereStr) > 0 {
		query += whereStr
		args = append(args, whereArg...)
	}

	query += b.buildGroupBy()

	havingStr, havingArg := b.buildHaving()
	if len(havingStr) > 0 {
		query += havingStr
		args = append(args, havingArg...)
	}

	query += b.buildOrderBy()
	query += b.buildOffset()
	query += b.buildLimit()
	return query, args
}
