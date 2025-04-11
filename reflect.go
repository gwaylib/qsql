package qsql

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	"github.com/gwaylib/errors"
	"github.com/jmoiron/sqlx/reflectx"
)

var refxM = reflectx.NewMapperTagFunc("db", func(in string) string {
	// for tag name
	return in
}, func(in string) string {
	// for options
	trims := []string{}
	options := strings.Split(in, ",")
	for _, op := range options {
		trims = append(trims, strings.TrimSpace(op))
	}
	return strings.Join(trims, ",")
})

// return is it a auto_increment field
func travelStructField(f *reflectx.FieldInfo, v *reflect.Value, order *int, drvName *string, outputNames *[]byte, outputInputs *[]byte, outputVals *[]interface{}) *reflect.Value {
	*order += 1
	switch v.Kind() {
	case reflect.Invalid:
		// nil value
		return nil
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		// continue
		break
	case reflect.Struct, reflect.Ptr:
		if _, ok := v.Interface().(driver.Valuer); ok {
			break
		}
		switch v.Type().String() {
		case "time.Time":
			break
		default:
			var autoIncrement *reflect.Value
			childrenLen := len(f.Children)
			for i := 0; i < childrenLen; i++ {
				child := f.Children[i]
				if child == nil {
					// found ignore tag, do next.
					continue
				}
				fieldVal := reflect.Indirect(*v).Field(i)
				autoFiled := travelStructField(
					child,
					&fieldVal,
					order, drvName,
					outputNames, outputInputs, outputVals,
				)
				if autoFiled != nil {
					autoIncrement = autoFiled
				}
			}
			return autoIncrement
		}
	default:
		// unsupport
		switch v.Type().String() {
		case "[]uint8":
			break
		default:
			return nil
		}
	}

	//
	// decode fileds
	//

	_, ok := f.Options["autoincrement"]
	if ok {
		// ignore 'autoincrement' for insert data
		return v
	}
	_, ok = f.Options["auto_increment"]
	if ok {
		// ignore 'auto_increment' for insert data
		return v
	}

	*outputVals = append(*outputVals, v.Interface())
	switch *drvName {
	case DRV_NAME_ORACLE, _DRV_NAME_OCI8:
		*order += 1
		*outputNames = append(*outputNames, []byte("\""+f.Name+"\",")...)
		*outputInputs = append(*outputInputs, []byte(fmt.Sprintf(":%s,", f.Name))...)
	case DRV_NAME_POSTGRES:
		*outputNames = append(*outputNames, []byte("\""+f.Name+"\",")...)
		*outputInputs = append(*outputInputs, []byte(fmt.Sprintf(":%d,", *order))...)
		*order += 1
	case DRV_NAME_SQLSERVER, _DRV_NAME_MSSQL:
		*outputNames = append(*outputNames, []byte("["+f.Name+"],")...)
		*outputInputs = append(*outputInputs, []byte(fmt.Sprintf("@p%d,", *order))...)
		*order += 1
	case DRV_NAME_MYSQL:
		*order += 1
		*outputNames = append(*outputNames, []byte("`"+f.Name+"`,")...)
		*outputInputs = append(*outputInputs, []byte("?,")...)
	default:
		*outputNames = append(*outputNames, []byte("\""+f.Name+"\",")...)
		*outputInputs = append(*outputInputs, []byte("?,")...)
	}

	return nil
}

type reflectInsertField struct {
	Names  string
	Stmts  string
	Values []interface{}

	AutoIncrement *reflect.Value
}

func (r *reflectInsertField) SetAutoIncrement(v reflect.Value) {
	if r.AutoIncrement == nil {
		return
	}
	r.AutoIncrement.Set(v)
}

func reflectInsertStruct(i interface{}, drvName string) (*reflectInsertField, error) {
	v := reflect.ValueOf(i)
	k := v.Kind()
	switch k {
	case reflect.Ptr:
	default:
		return nil, errors.New("Unsupport reflect type").As(k.String())
	}
	v = reflect.Indirect(v)

	tm := refxM.TypeMap(v.Type())

	names := []byte{}
	inputs := []byte{}
	vals := []interface{}{}
	var autoIncrement *reflect.Value

	childrenLen := len(tm.Tree.Children)
	order := 0
	for i := 0; i < childrenLen; i++ {
		field := tm.Tree.Children[i]
		if field == nil {
			// found ignore tag, do next.
			continue
		}

		fieldVal := v.Field(i)
		autoField := travelStructField(field, &fieldVal, &order, &drvName, &names, &inputs, &vals)
		if autoField != nil {
			autoIncrement = autoField
		}
	}

	if len(names) == 0 {
		panic("No public field in struct")
	}
	return &reflectInsertField{
		Names:         string(names[:len(names)-1]),
		Stmts:         string(inputs[:len(inputs)-1]),
		Values:        vals,
		AutoIncrement: autoIncrement,
	}, nil
}
