package qsql

import (
	"fmt"
)

func stmtIn(paramIdx, paramsLen int, driverNames ...string) string {
	driverName := ""
	if len(driverNames) > 0 {
		driverName = driverNames[0]
	}
	drvName := getDrvName(nil, driverName)
	switch drvName {
	case DRV_NAME_ORACLE, _DRV_NAME_OCI8:
		result := []byte{}
		for i := 0; i < paramsLen; i++ {
			result = append(result, []byte(fmt.Sprintf(":%d,", paramIdx+i))...)
		}
		if len(result) > 0 {
			return string(result[:len(result)-1]) // remove the last ','
		}
		return string(result)
	case DRV_NAME_POSTGRES:
		result := []byte{}
		for i := 0; i < paramsLen; i++ {
			result = append(result, []byte(fmt.Sprintf("$%d,", paramIdx+i))...)
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
