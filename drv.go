package qsql

import "github.com/gwaylib/errors"

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
