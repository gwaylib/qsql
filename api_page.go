// Example:
//
// mdb := db.GetCache("main")
//
// // count sql
// cbd := NewSqlBuilder(mdb.DriverName())
// cbd.Select("COUNT(*)")
// cbd.Add("FROM tmp")
// cbd.Add("WHERE")
// cbd.AddTab("create_at BETWEEN ? AND ?", time.Now().AddDate(-1,0,0), time.Now())
//
// // copy condition
// qbd := cbd.Copy()
// qbd.Select("id", "created_at", "name")
// qbd.Add("OFFSET ?", 0)
// qbd.Add("LIMIT ?", 20)
//
// pSql := NewPageSql(cbd, qbd)
// count, err := pSql.QueryCount(db)
// ...
// Or
// titles, result, err := pSql.QueryPageArray(db)
// ...
// Or
// titles, result, err := pSql.QueryPageMap(db)
// ...
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
