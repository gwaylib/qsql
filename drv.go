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

func getDrvName(exec Execer, driverName ...string) string {
	drvName := ""
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
	if len(drvName) == 0 {
		panic("driver name not set")
	}
	return drvName
}
