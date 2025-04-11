package qsql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

// Bool type
type Bool bool

func (v *Bool) Scan(i interface{}) error {
	b := sql.NullBool{}
	if err := b.Scan(i); err != nil {
		return err
	}
	*v = Bool(b.Bool)
	return nil
}
func (v *Bool) Value() (driver.Value, error) {
	return v, nil
}

// Int64 type
type Int64 int64

func (v *Int64) Scan(i interface{}) error {
	b := sql.NullInt64{}
	if err := b.Scan(i); err != nil {
		return err
	}
	*v = Int64(b.Int64)
	return nil
}
func (v Int64) Value() (driver.Value, error) {
	return int64(v), nil
}

// Float64 type
type Float64 float64

func (v *Float64) Scan(i interface{}) error {
	b := sql.NullFloat64{}
	if err := b.Scan(i); err != nil {
		return err
	}
	*v = Float64(b.Float64)
	return nil
}
func (v Float64) Value() (driver.Value, error) {
	return float64(v), nil
}

// String type
type String string

func (v *String) Scan(i interface{}) error {
	b := sql.NullString{}
	if err := b.Scan(i); err != nil {
		return err
	}
	*v = String(b.String)
	return nil
}
func (v String) Value() (driver.Value, error) {
	return string(v), nil
}
func (v *String) String() string {
	return string(*v)
}

// 通用的字符串查询
type DBData string

func (d *DBData) Scan(i interface{}) error {
	if i == nil {
		*d = ""
		return nil
	}
	switch i.(type) {
	case int64:
		*d = DBData(strconv.FormatInt(i.(int64), 10))
	case float64:
		*d = DBData(strconv.FormatFloat(i.(float64), 'f', -1, 64))
	case []byte:
		*d = DBData(string(i.([]byte)))
	case string:
		*d = DBData(i.(string))
	case bool:
		*d = DBData(fmt.Sprintf("%t", i))
	case time.Time:
		*d = DBData(i.(time.Time).Format(time.RFC3339))
	default:
		*d = DBData(fmt.Sprint(i))
	}
	return nil
}
func (d *DBData) String() string {
	return string(*d)
}

func MakeDBData(l int) []interface{} {
	r := make([]interface{}, l)
	for i := 0; i < l; i++ {
		d := DBData("")
		r[i] = &d
	}
	return r
}
