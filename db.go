package qsql

import (
	"context"
	"database/sql"
	"sync"
)

// qsql.DB Extendd sql.DB
// and implement qsql.QuickQuery interface
type DB struct {
	*sql.DB
	drvName string
	isClose bool
	mu      sync.Mutex
}

func _checkQuickSql() QuickSql {
	return &DB{}
}

func newDB(drvName string, db *sql.DB) *DB {
	return &DB{
		DB:      db,
		drvName: drvName,
	}
}

func (db *DB) DriverName() string {
	return db.drvName
}

func (db *DB) IsClose() bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.isClose
}

func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.isClose = true
	rmCache(db)
	return db.DB.Close()
}

func (db *DB) NewSqlBuilder() *SqlBuilder {
	return NewSqlBuilder(db.drvName)
}

// Reflect one db data to the struct. the struct tag format like `db:"field_title"`, reference to: http://github.com/jmoiron/sqlx
func (db *DB) InsertStruct(structPtr interface{}, tbName string) (sql.Result, error) {
	return insertStruct(db, context.TODO(), structPtr, tbName, db.drvName)
}
func (db *DB) InsertStructContext(ctx context.Context, structPtr interface{}, tbName string) (sql.Result, error) {
	return insertStruct(db, ctx, structPtr, tbName, db.drvName)
}

// Reflect the sql.Rows to []struct array.
// Return empty array if data not found.
// Refere to: github.com/jmoiron/sqlx
// DO NOT forget close the rows
func (db *DB) ScanStructs(rows *sql.Rows, structsPtr interface{}) error {
	return scanStructs(rows, structsPtr)
}

// Reflect the sql.Query result to a struct.
func (db *DB) QueryStruct(structPtr interface{}, querySql string, args ...interface{}) error {
	return queryStruct(db, context.TODO(), structPtr, querySql, args...)
}
func (db *DB) QueryStructContext(ctx context.Context, structPtr interface{}, querySql string, args ...interface{}) error {
	return queryStruct(db, ctx, structPtr, querySql, args...)
}

// Reflect the sql.Query result to a struct array.
// Return empty array if data not found.
func (db *DB) QueryStructs(structPtr interface{}, querySql string, args ...interface{}) error {
	return queryStructs(db, context.TODO(), structPtr, querySql, args...)
}
func (db *DB) QueryStructsContext(ctx context.Context, structPtr interface{}, querySql string, args ...interface{}) error {
	return queryStructs(db, ctx, structPtr, querySql, args...)
}

// Query one field to a sql.Scanner.
func (db *DB) QueryElem(result interface{}, querySql string, args ...interface{}) error {
	return queryElem(db, context.TODO(), result, querySql, args...)
}
func (db *DB) QueryElemContext(ctx context.Context, result interface{}, querySql string, args ...interface{}) error {
	return queryElem(db, ctx, result, querySql, args...)
}

// Query one field to a sql.Scanner array.
func (db *DB) QueryElems(result interface{}, querySql string, args ...interface{}) error {
	return queryElems(db, context.TODO(), result, querySql, args...)
}
func (db *DB) QueryElemsContext(ctx context.Context, result interface{}, querySql string, args ...interface{}) error {
	return queryElems(db, ctx, result, querySql, args...)
}

// Reflect the query result to a string array.
func (db *DB) QueryPageArr(querySql string, args ...interface{}) (titles []string, result [][]interface{}, err error) {
	return queryPageArr(db, context.TODO(), querySql, args...)
}
func (db *DB) QueryPageArrContext(ctx context.Context, querySql string, args ...interface{}) (titles []string, result [][]interface{}, err error) {
	return queryPageArr(db, ctx, querySql, args...)
}

// Reflect the query result to a string map.
func (db *DB) QueryPageMap(querySql string, args ...interface{}) (titles []string, result []map[string]interface{}, err error) {
	return queryPageMap(db, context.TODO(), querySql, args...)
}
func (db *DB) QueryPageMapContext(ctx context.Context, querySql string, args ...interface{}) (titles []string, result []map[string]interface{}, err error) {
	return queryPageMap(db, ctx, querySql, args...)
}

// Return "?,?,?,?..." for default, or "@p1,@p2,@p3..." for mssql, or ":1,:2,:3..." for pgsql.
// paramStartIdx default is 0, but you need count it when the driver is mssq, pgsql etc. .
func (db *DB) StmtIn(paramStartIdx, paramsLen int) string {
	return stmtIn(paramStartIdx, paramsLen, db.DriverName())
}

// A lazy function to commit the *sql.Tx
// if will auto commit when the function is nil error, or do a rollback and return the function error.
func (db *DB) Commit(tx *sql.Tx, fn func() error) error {
	if err := fn(); err != nil {
		Rollback(tx)
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
