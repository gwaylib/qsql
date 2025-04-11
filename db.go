package qsql

import (
	"context"
	"database/sql"
	"sync"
)

type DB struct {
	*sql.DB
	drvName string
	isClose bool
	mu      sync.Mutex
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

// Reflect one db data to the struct. the struct tag format like `db:"field_title"`, reference to: http://github.com/jmoiron/sqlx
func (db *DB) InsertStruct(obj interface{}, tbName string) (sql.Result, error) {
	return insertStruct(db, context.TODO(), obj, tbName, db.drvName)
}
func (db *DB) InsertStructContext(ctx context.Context, obj interface{}, tbName string) (sql.Result, error) {
	return insertStruct(db, ctx, obj, tbName, db.drvName)
}

// Relect the sql.Rows to a struct.
func (db *DB) ScanStruct(rows Rows, obj interface{}) error {
	return scanStruct(rows, obj)
}

// Reflect the sql.Rows to a struct array.
// Return empty array if data not found.
// Refere to: github.com/jmoiron/sqlx
func (db *DB) ScanStructs(rows Rows, obj interface{}) error {
	return scanStructs(rows, obj)
}

// Reflect the sql.Query result to a struct.
func (db *DB) QueryStruct(obj interface{}, querySql string, args ...interface{}) error {
	return queryStruct(db, context.TODO(), obj, querySql, args...)
}
func (db *DB) QueryStructContext(ctx context.Context, obj interface{}, querySql string, args ...interface{}) error {
	return queryStruct(db, ctx, obj, querySql, args...)
}

// Reflect the sql.Query result to a struct array.
// Return empty array if data not found.
func (db *DB) QueryStructs(obj interface{}, querySql string, args ...interface{}) error {
	return queryStructs(db, context.TODO(), obj, querySql, args...)
}
func (db *DB) QueryStructsContext(ctx context.Context, obj interface{}, querySql string, args ...interface{}) error {
	return queryStructs(db, ctx, obj, querySql, args...)
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

// Return "?,?,?,?..." for default, or "@p1,@p2,@p3..." for mssql, or ":1,:2,:3..." for pgsql when paramStartIdx is 0.
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
