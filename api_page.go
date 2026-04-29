package qsql

import (
	"context"
	"fmt"

	"github.com/gwaylib/errors"
)

type PageSql struct {
	countSql string
	querySql string
}

func NewPageSql(countSql, querySql string) *PageSql {
	return &PageSql{
		countSql: countSql,
		querySql: querySql,
	}
}

// format sql and return a new PageSql
func (p *PageSql) FmtPage(args ...interface{}) *PageSql {
	return &PageSql{
		countSql: fmt.Sprintf(p.countSql, args...),
		querySql: fmt.Sprintf(p.querySql, args...),
	}
}

func (p *PageSql) QueryCount(db *DB, args ...interface{}) (int64, error) {
	count := int64(0)
	if err := queryElem(db, context.TODO(), &count, p.countSql, args...); err != nil {
		return 0, errors.As(err)
	}
	return count, nil
}

func (p *PageSql) QueryPageArr(db *DB, args ...interface{}) ([]string, [][]interface{}, error) {
	titles, data, err := queryPageArr(db, context.TODO(), p.querySql, args...)
	if err != nil {
		return nil, nil, errors.As(err)
	}
	return titles, data, nil
}

func (p *PageSql) QueryPageMap(db *DB, args ...interface{}) ([]string, []map[string]interface{}, error) {
	titles, data, err := queryPageMap(db, context.TODO(), p.querySql, args...)
	if err != nil {
		return nil, nil, errors.As(err)
	}
	return titles, data, nil
}
