package qsql

import (
	"context"

	"github.com/gwaylib/errors"
)

type PageSql struct {
	countBD *SqlBuilder
	queryBD *SqlBuilder
}

func NewPageSql(countBD, queryBD *SqlBuilder) *PageSql {
	return &PageSql{
		countBD: countBD,
		queryBD: queryBD,
	}
}

func (p *PageSql) QueryCount(db *DB) (int64, error) {
	count := int64(0)
	if err := queryElem(db, context.TODO(), &count, p.countBD.String(), p.countBD.Args()...); err != nil {
		return 0, errors.As(err)
	}
	return count, nil
}

func (p *PageSql) QueryPageArr(db *DB) ([]string, [][]interface{}, error) {
	titles, data, err := queryPageArr(db, context.TODO(), p.queryBD.String(), p.queryBD.Args()...)
	if err != nil {
		return nil, nil, errors.As(err)
	}
	return titles, data, nil
}

func (p *PageSql) QueryPageMap(db *DB) ([]string, []map[string]interface{}, error) {
	titles, data, err := queryPageMap(db, context.TODO(), p.queryBD.String(), p.queryBD.Args()...)
	if err != nil {
		return nil, nil, errors.As(err)
	}
	return titles, data, nil
}
