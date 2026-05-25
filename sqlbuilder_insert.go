package qsql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/gwaylib/errors"
)

const (
	addObjSql = "INSERT INTO %s (%s) VALUES (%s);"
)

// field flag like: `db:"name"`
// more: github.com/jmoiron/sqlx
func insertStruct(exec Execer, ctx context.Context, obj interface{}, tbName string, driverName ...string) (sql.Result, error) {
	drvName := getDrvName(exec, driverName...)

	fields, err := reflectInsertStruct(obj, drvName)
	if err != nil {
		return nil, errors.As(err)
	}
	execSql := fmt.Sprintf(addObjSql, tbName, strings.Join(fields.Names, ", "), strings.Join(fields.Stmts, ", "))
	// log.Debugf("%s%+v", execSql, vals)
	result, err := exec.ExecContext(ctx, execSql, fields.Values...)
	if err != nil {
		return nil, errors.As(err, execSql)
	}
	if fields.AutoIncrement != nil {
		id, _ := result.LastInsertId()
		var val reflect.Value
		kind := fields.AutoIncrement.Kind()
		switch kind {
		case reflect.Int:
			val = reflect.ValueOf(int(id))
		case reflect.Int8:
			val = reflect.ValueOf(int8(id))
		case reflect.Int16:
			val = reflect.ValueOf(int16(id))
		case reflect.Int32:
			val = reflect.ValueOf(int32(id))
		case reflect.Int64:
			val = reflect.ValueOf(int64(id))
		case reflect.Uint: // Warnning: this maybe out of int64
			val = reflect.ValueOf(uint(id))
		case reflect.Uint8:
			val = reflect.ValueOf(uint8(id))
		case reflect.Uint16:
			val = reflect.ValueOf(uint16(id))
		case reflect.Uint32:
			val = reflect.ValueOf(uint32(id))
		case reflect.Uint64: // Warnning: this maybe out of int64
			val = reflect.ValueOf(uint64(id))
		default:
			// unsupport other kind here
			panic("unsupport auto increment kind: " + kind.String())
		}

		fields.AutoIncrement.Set(val)
	}
	return result, nil
}
