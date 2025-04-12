package qsql

import (
	"fmt"
	"strconv"
	"time"
)

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
