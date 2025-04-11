package qsql

type SqlBuilder struct {
	table          []string
	selectStr      []string
	whereAndStr    []string
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

func (b *SqlBuilder) Joins(table string) *SqlBuilder {
	b.table = append(b.table, table)
	return b
}

func (b *SqlBuilder) Where(cond string, args ...interface{}) *SqlBuilder {
	b.whereAndStr = append(b.whereAndStr, cond)
	b.whereParams = append(b.whereParams, args...)
	return b
}
