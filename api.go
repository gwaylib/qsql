/*
Provides database connections in factory mode to optimize database connections
*/
package qsql

import (
	"context"
	"database/sql"
	"io"
	"runtime/debug"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
)

type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row

	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type QuickSql interface {
	DriverName() string

	// Insert a struct data into tbName
	//
	// Reflect one db data to the struct.
	// the struct tag format like `db:"field_title"`, reference to: http://github.com/jmoiron/sqlx
	//
	InsertStruct(structPtr interface{}, tbName string) (sql.Result, error)
	InsertStructContext(ctx context.Context, structPtr interface{}, tbName string) (sql.Result, error)

	// Scan the rows result to []struct
	// Reflect the sql.Rows to a struct array.
	// Return empty array if data not found.
	// Refere to: github.com/jmoiron/sqlx
	// DOT NOT forget close the rows after called.
	ScanStructs(rows *sql.Rows, structsPtr interface{}) error

	// Query db data to a struct
	QueryStruct(structPrt interface{}, querySql string, args ...interface{}) error
	QueryStructContext(ctx context.Context, structPrt interface{}, querySql string, args ...interface{}) error
	// Query db data to []struct
	QueryStructs(structsPrt interface{}, querySql string, args ...interface{}) error
	QueryStructsContext(ctx context.Context, structsPrt interface{}, querySql string, args ...interface{}) error

	// Query a element data like int, string.
	// Same as row.Scan(&e)
	QueryElem(ePtr interface{}, querySql string, args ...interface{}) error
	QueryElemContext(ctx context.Context, ePtr interface{}, querySql string, args ...interface{}) error
	// Query elements data like []int, []string in result.
	QueryElems(ePtr interface{}, querySql string, args ...interface{}) error
	QueryElemsContext(ctx context.Context, ePtr interface{}, querySql string, args ...interface{}) error

	// Query a page data to array.
	// the result data is [][]*string but no nil *string pointer instance.
	QueryPageArr(querySql string, args ...interface{}) (titles []string, result [][]interface{}, err error)
	QueryPageArrContext(ctx context.Context, querySql string, args ...interface{}) (titles []string, result [][]interface{}, err error)

	// Query a page data to map, NOT RECOMMENED to use when there is a large page data.
	// the result data is []map[string]*string but no nil *string pointer instance.
	QueryPageMap(querySql string, args ...interface{}) (titles []string, result []map[string]interface{}, err error)
	QueryPageMapContext(ctx context.Context, querySql string, args ...interface{}) (titles []string, result []map[string]interface{}, err error)

	// Extend stmt for the where in
	// paramStartIdx default is 0, but you need count it when the driver is mssq, pgsql etc. .
	//
	// Example for the first input:
	// fmt.Sprintf("select * from table_name where in (%s)", qsql.StmtWhereIn(0,len(args))
	// Or
	// fmt.Sprintf("select * from table_name where in (%s)", qsql.StmtWhereIn(0,len(args), qsql.DRV_NAME_MYSQL)
	//
	// Example for the second input:
	// fmt.Sprintf("select * from table_name where id=? in (%s)", qsql.StmtWhereIn(1,len(args))
	//
	// Return "?,?,?,?..." for default, or "@p1,@p2,@p3..." for mssql, or ":1,:2,:3..." for pgsql when paramStartIdx is 0.
	MakeStmtIn(paramStartIdx, paramLen int) string

	// auto commit when the func is return nil, or auto rollback when the func is error
	Commit(tx *sql.Tx, fn func() error) error
}

func NewDB(drvName string, db *sql.DB) *DB {
	return newDB(drvName, db)
}

// Implement the sql.Open
func Open(drvName, dsn string) (*DB, error) {
	db, err := sql.Open(drvName, dsn)
	if err != nil {
		return nil, errors.As(err, drvName, dsn)
	}
	return newDB(drvName, db), nil
}

// A lazy function to closed the io.Closer
func Close(closer io.Closer) {
	if closer == nil {
		return
	}
	if err := closer.Close(); err != nil {
		println(errors.As(err).Error())
		debug.PrintStack()
	}
}

// A lazy function to rollback the *sql.Tx
func Rollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		// roll back error is a serious error
		log.Error(errors.As(err))
	}
}

// A lazy function to commit the *sql.Tx
func Commit(tx *sql.Tx, fn func() error) error {
	if err := fn(); err != nil {
		Rollback(tx)
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func InsertStruct(drvName string, exec Execer, obj interface{}, tbName string) (sql.Result, error) {
	return insertStruct(exec, context.TODO(), obj, tbName, drvName)
}
func InsertStructContext(drvName string, exec Execer, ctx context.Context, obj interface{}, tbName string) (sql.Result, error) {
	return insertStruct(exec, ctx, obj, tbName, drvName)
}

func ScanStructs(rows *sql.Rows, obj interface{}) error {
	return scanStructs(rows, obj)
}

func QueryStruct(queryer Queryer, obj interface{}, querySql string, args ...interface{}) error {
	return queryStruct(queryer, context.TODO(), obj, querySql, args...)
}
func QueryStructContext(queryer Queryer, ctx context.Context, obj interface{}, querySql string, args ...interface{}) error {
	return queryStruct(queryer, ctx, obj, querySql, args...)
}

func QueryStructs(queryer Queryer, obj interface{}, querySql string, args ...interface{}) error {
	return queryStructs(queryer, context.TODO(), obj, querySql, args...)
}
func QueryStructsContext(queryer Queryer, ctx context.Context, obj interface{}, querySql string, args ...interface{}) error {
	return queryStructs(queryer, ctx, obj, querySql, args...)
}

func QueryElem(queryer Queryer, result interface{}, querySql string, args ...interface{}) error {
	return queryElem(queryer, context.TODO(), result, querySql, args...)
}
func QueryElemContext(queryer Queryer, ctx context.Context, result interface{}, querySql string, args ...interface{}) error {
	return queryElem(queryer, ctx, result, querySql, args...)
}

func QueryElems(queryer Queryer, result interface{}, querySql string, args ...interface{}) error {
	return queryElems(queryer, context.TODO(), result, querySql, args...)
}
func QueryElemsContext(queryer Queryer, ctx context.Context, result interface{}, querySql string, args ...interface{}) error {
	return queryElems(queryer, ctx, result, querySql, args...)
}

func QueryPageArr(queryer Queryer, querySql string, args ...interface{}) (titles []string, result [][]interface{}, err error) {
	return queryPageArr(queryer, context.TODO(), querySql, args...)
}
func QueryPageArrContext(queryer Queryer, ctx context.Context, querySql string, args ...interface{}) (titles []string, result [][]interface{}, err error) {
	return queryPageArr(queryer, ctx, querySql, args...)
}

func QueryPageMap(queryer Queryer, querySql string, args ...interface{}) (titles []string, result []map[string]interface{}, err error) {
	return queryPageMap(queryer, context.TODO(), querySql, args...)
}
func QueryPageMapContext(queryer Queryer, ctx context.Context, querySql string, args ...interface{}) (titles []string, result []map[string]interface{}, err error) {
	return queryPageMap(queryer, ctx, querySql, args...)
}

func StmtIn(paramStartIdx, paramsLen int, drvName ...string) string {
	return stmtIn(paramStartIdx, paramsLen, drvName...)
}
