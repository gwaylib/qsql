package qsql

import (
	"fmt"
	"reflect"
)

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

func StmtSliceArgs(args ...interface{}) []interface{} {
	result := []interface{}{}
	for _, arg := range args {
		val := reflect.ValueOf(arg)
		switch val.Kind() {
		case reflect.Array, reflect.Slice:
			arrLen := val.Len()
			for i := 0; i < arrLen; i++ {
				result = append(result, val.Index(i).Interface())
			}
		default:
			result = append(result, arg)
		}
	}
	return result
}
