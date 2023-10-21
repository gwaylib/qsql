/*
Provides database connections in factory mode to optimize database connections
*/
package qsql

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"runtime/debug"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
)

const (
	DRV_NAME_MYSQL     = "mysql"
	DRV_NAME_ORACLE    = "oracle" // or "oci8"
	DRV_NAME_POSTGRES  = "postgres"
	DRV_NAME_SQLITE3   = "sqlite3"
	DRV_NAME_SQLSERVER = "sqlserver" // or "mssql"

	_DRV_NAME_OCI8  = "oci8"
	_DRV_NAME_MSSQL = "mssql"
)

var (
	// Whe reflect the QueryStruct, InsertStruct, it need set the Driver first.
	// For example:
	// func init(){
	//     qsql.REFLECT_DRV_NAME = qsql.DEV_NAME_SQLITE3
	// }
	// Default is using the mysql driver.
	REFLECT_DRV_NAME = DRV_NAME_MYSQL
)

func getDrvName(exec Execer, driverName ...string) string {
	drvName := REFLECT_DRV_NAME
	db, ok := exec.(*DB)
	if ok {
		drvName = db.DriverName()
	} else {
		drvNamesLen := len(driverName)
		if drvNamesLen > 0 {
			if drvNamesLen != 1 {
				panic(errors.New("'drvName' expect only one argument").As(driverName))
			}
			drvName = driverName[0]
		}
	}
	return drvName
}

// Extend the where in stmt
//
// Example for the first input:
// fmt.Sprintf("select * from table_name where in (%s)", qsql.StmtWhereIn(0,len(args))
// Or
// fmt.Sprintf("select * from table_name where in (%s)", qsql.StmtWhereIn(0,len(args), qsql.DRV_NAME_MYSQL)
//
// Example for the second input:
// fmt.Sprintf("select * from table_name where id=? in (%s)", qsql.StmtWhereIn(1,len(args))
//
func StmtWhereIn(paramIdx, paramsLen int, driverName ...string) string {
	drvName := getDrvName(nil, driverName...)
	switch drvName {
	case DRV_NAME_ORACLE, _DRV_NAME_OCI8:
		// *outputInputs = append(*outputInputs, []byte(fmt.Sprintf(":%s,", f.Name))...)
		panic("unknow how to implemented")
	case DRV_NAME_POSTGRES:
		result := []byte{}
		for i := 0; i < paramsLen; i++ {
			result = append(result, []byte(fmt.Sprintf(":%d,", paramIdx+i))...)
		}
		if len(result) > 0 {
			return string(result[:len(result)-1]) // remove the last ','
		}
		return string(result)
	case DRV_NAME_SQLSERVER, _DRV_NAME_MSSQL:
		result := []byte{}
		for i := 0; i < paramsLen; i++ {
			result = append(result, []byte(fmt.Sprintf("@p%d,", paramIdx+i))...)
		}
		if len(result) > 0 {
			return string(result[:len(result)-1]) // remove the last ','
		}
		return string(result)
	default:
		resultLen := paramsLen * 2
		result := make([]byte, resultLen)
		for i := 0; i < resultLen; i += 2 {
			result[i] = '?'
			result[i+1] = ','
		}
		if len(result) > 0 {
			return string(result[:len(result)-1]) // remove the last ','
		}
		return string(result)
	}
}

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

type Rows interface {
	Close() error
	Columns() ([]string, error)
	Err() error
	Next() bool
	Scan(...interface{}) error
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

// Register a db to the connection pool by manully.
func RegCache(iniFileName, sectionName string, db *DB) {
	regCache(iniFileName, sectionName, db)
}

// Get the db instance from the cache.
// If the db not in the cache, it will create a new instance from the ini file.
func GetCache(iniFileName, sectionName string) *DB {
	db, err := getCache(iniFileName, sectionName)
	if err != nil {
		panic(err)
	}
	return db
}

// Checking the cache does it have a db instance.
func HasCache(etcFileName, sectionName string) (*DB, error) {
	return getCache(etcFileName, sectionName)
}

// Close all instance in the cache.
func CloseCache() {
	closeCache()
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
	err := tx.Rollback()

	// roll back error is a serious error
	if err != nil {
		log.Error(errors.As(err))
	}
}

// A way implement the sql.Exec
func Exec(db Execer, querySql string, args ...interface{}) (sql.Result, error) {
	return db.Exec(querySql, args...)
}
func ExecContext(db Execer, ctx context.Context, querySql string, args ...interface{}) (sql.Result, error) {
	return db.ExecContext(ctx, querySql, args...)
}

// A way to ran multiply tx
func ExecMultiTx(tx *sql.Tx, mTx []*MultiTx) error {
	return execMultiTx(tx, context.TODO(), mTx)
}
func ExecMultiTxContext(tx *sql.Tx, ctx context.Context, mTx []*MultiTx) error {
	return execMultiTx(tx, ctx, mTx)
}

// Reflect one db data to the struct. the struct tag format like `db:"field_title"`, reference to: http://github.com/jmoiron/sqlx
// When you no set the REFLECT_DRV_NAME, you can point out with the drvName
func InsertStruct(exec Execer, obj interface{}, tbName string, drvName ...string) (sql.Result, error) {
	return insertStruct(exec, context.TODO(), obj, tbName, drvName...)
}
func InsertStructContext(exec Execer, ctx context.Context, obj interface{}, tbName string, drvName ...string) (sql.Result, error) {
	return insertStruct(exec, ctx, obj, tbName, drvName...)
}

// A sql.Query implements
func Query(db Queryer, querySql string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(querySql, args...)
}
func QueryContext(db Queryer, ctx context.Context, querySql string, args ...interface{}) (*sql.Rows, error) {
	return db.QueryContext(ctx, querySql, args...)
}

// A sql.QueryRow implements
func QueryRow(db Queryer, querySql string, args ...interface{}) *sql.Row {
	return db.QueryRow(querySql, args...)
}
func QueryRowContext(db Queryer, ctx context.Context, querySql string, args ...interface{}) *sql.Row {
	return db.QueryRowContext(ctx, querySql, args...)
}

// Relect the sql.Rows to a struct.
func ScanStruct(rows Rows, obj interface{}) error {
	return scanStruct(rows, obj)
}

// Reflect the sql.Rows to a struct array.
// Return empty array if data not found.
// Refere to: github.com/jmoiron/sqlx
func ScanStructs(rows Rows, obj interface{}) error {
	return scanStructs(rows, obj)
}

// Reflect the sql.Query result to a struct.
func QueryStruct(db Queryer, obj interface{}, querySql string, args ...interface{}) error {
	return queryStruct(db, context.TODO(), obj, querySql, args...)
}
func QueryStructContext(db Queryer, ctx context.Context, obj interface{}, querySql string, args ...interface{}) error {
	return queryStruct(db, ctx, obj, querySql, args...)
}

// Reflect the sql.Query result to a struct array.
// Return empty array if data not found.
func QueryStructs(db Queryer, obj interface{}, querySql string, args ...interface{}) error {
	return queryStructs(db, context.TODO(), obj, querySql, args...)
}
func QueryStructsContext(db Queryer, ctx context.Context, obj interface{}, querySql string, args ...interface{}) error {
	return queryStructs(db, ctx, obj, querySql, args...)
}

// Query one field to a sql.Scanner.
func QueryElem(db Queryer, result interface{}, querySql string, args ...interface{}) error {
	return queryElem(db, context.TODO(), result, querySql, args...)
}
func QueryElemContext(db Queryer, ctx context.Context, result interface{}, querySql string, args ...interface{}) error {
	return queryElem(db, ctx, result, querySql, args...)
}

// Query one field to a sql.Scanner array.
func QueryElems(db Queryer, result interface{}, querySql string, args ...interface{}) error {
	return queryElems(db, context.TODO(), result, querySql, args...)
}
func QueryElemsContext(db Queryer, ctx context.Context, result interface{}, querySql string, args ...interface{}) error {
	return queryElems(db, ctx, result, querySql, args...)
}

// Reflect the query result to a string array.
func QueryPageArr(db Queryer, querySql string, args ...interface{}) (titles []string, result [][]interface{}, err error) {
	return queryPageArr(db, context.TODO(), querySql, args...)
}
func QueryPageArrContext(db Queryer, ctx context.Context, querySql string, args ...interface{}) (titles []string, result [][]interface{}, err error) {
	return queryPageArr(db, ctx, querySql, args...)
}

// Reflect the query result to a string map.
func QueryPageMap(db Queryer, querySql string, args ...interface{}) (titles []string, result []map[string]interface{}, err error) {
	return queryPageMap(db, context.TODO(), querySql, args...)
}
func QueryPageMapContext(db Queryer, ctx context.Context, querySql string, args ...interface{}) (titles []string, result []map[string]interface{}, err error) {
	return queryPageMap(db, ctx, querySql, args...)
}
