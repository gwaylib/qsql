package qsql

type SqlBuilder struct {
	innerBuilder *SqlBuilder
	table        []string

	whereStr       []string
	whereParams    []interface{}
	whereOrStr     []string
	whereOrParams  []interface{}
	groupByStr     []string
	havingStr      []string
	havingParams   []interface{}
	havingOrStr    []string
	havingOrParams []interface{}
	orderByStr     []string
	offset         int
	limit          int
}

func NewSqlBuilder(table string) *SqlBuilder {
	return &SqlBuilder{
		table: []string{table},
	}
}

func (b *SqlBuilder) Copy() *SqlBuilder {
	n := &SqlBuilder{
		innerBuilder:   b.innerBuilder,
		table:          make([]string, len(b.table)),
		whereStr:       make([]string, len(b.whereStr)),
		whereParams:    make([]interface{}, len(b.whereParams)),
		whereOrStr:     make([]string, len(b.whereOrStr)),
		whereOrParams:  make([]interface{}, len(b.whereOrParams)),
		groupByStr:     make([]string, len(b.groupByStr)),
		havingStr:      make([]string, len(b.havingStr)),
		havingParams:   make([]interface{}, len(b.havingParams)),
		havingOrStr:    make([]string, len(b.havingOrStr)),
		havingOrParams: make([]interface{}, len(b.havingOrParams)),
		orderByStr:     make([]string, len(b.orderByStr)),
		offset:         b.offset,
		limit:          b.limit,
	}
	copy(n.table, b.table)
	copy(n.whereStr, b.whereStr)
	copy(n.whereParams, b.whereParams)
	copy(n.whereOrStr, b.whereOrStr)
	copy(n.whereOrParams, b.whereOrParams)
	copy(n.groupByStr, b.groupByStr)
	copy(n.havingStr, b.havingStr)
	copy(n.havingParams, b.havingParams)
	copy(n.havingOrStr, b.havingOrStr)
	copy(n.havingOrParams, b.havingOrParams)

	return n
}

func (b *SqlBuilder) Joins(table string) *SqlBuilder {
	b.table = append(b.table, table)
	return b
}

func (b *SqlBuilder) Where(cond string, args ...interface{}) *SqlBuilder {
	b.whereStr = append(b.whereStr, cond)
	b.whereParams = append(b.whereParams, args...)
	return b
}

func (b *SqlBuilder) Or(cond string, args ...interface{}) *SqlBuilder {
	b.whereOrStr = append(b.whereOrStr, cond)
	b.whereOrParams = append(b.whereOrParams, args...)
	return b
}

func (b *SqlBuilder) GroupBy(cond ...string) *SqlBuilder {
	b.groupByStr = append(b.groupByStr, cond...)
	return b
}
func (b *SqlBuilder) Having(cond string, args ...interface{}) *SqlBuilder {
	b.havingStr = append(b.havingStr, cond)
	b.havingParams = append(b.havingParams, args...)
	return b
}
func (b *SqlBuilder) HavingOr(cond string, args ...interface{}) *SqlBuilder {
	b.havingOrStr = append(b.havingOrStr, cond)
	b.havingOrParams = append(b.havingOrParams, args...)
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
func (b *SqlBuilder) Select(column ...string) (string, []interface{}) {
	return "TODO", nil
}
