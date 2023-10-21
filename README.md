# Refere to:
```
database/sql
https://github.com/jmoiron/sqlx
```

# Example:
More example see the [example](./example) directory.

## Using etc cache
Assume that the configuration file path is: './etc/db.cfg'

The etc file content
```
[master]
driver: mysql
dsn: username:passwd@tcp(127.0.0.1:3306)/main?timeout=30s&strict=true&loc=Local&parseTime=true&allowOldPasswords=1
life_time:7200

[log]
driver: mysql
dsn: username:passwd@tcp(127.0.0.1:3306)/log?timeout=30s&strict=true&loc=Local&parseTime=true&allowOldPasswords=1
life_time:7200
```

Make a package for connection cache
``` text
package db

import (
	"github.com/gwaylib/conf"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/qsql"
	_ "github.com/go-sql-driver/mysql"
)

var dbFile = conf.RootDir() + "/etc/db.cfg"

func init() {
   qsql.REFLECT_DRV_NAME = qsql.DRV_NAME_MYSQL 
}

func GetCache(section string) *qsql.DB {
	return qsql.GetCache(dbFile, section)
}

func HasCache(section string) (*qsql.DB, error) {
	return qsql.HasCache(dbFile, section)
}

func CloseCache() {
	qsql.CloseCache()
}
```

Call a cache
``` text
mdb := db.GetCache("master")
```

## Call standar sql
``` text
mdb := db.GetCache("master") 
// or mdb = <sql.Tx>

// row := mdb.QueryRow("SELECT * ...")
row := qsql.QueryRow(mdb, "SELECT * ...")
// ...

// rows, err := mdb.Query("SELECT * ...")
rows, err := qsql.Query(mdb, "SELECT * ...")
// ...

// result, err := mdb.Exec("UPDATE ...")
result, err := qsql.Exec(mdb, "UPDATE ...")
// ...
```

## Insert a struct to db(using reflect)
``` text
type User struct{
    Id     int64  `db:"id,auto_increment"` // flag "autoincrement", "auto_increment" are supported .
    Name   string `db:"name"`
    Ignore string `db:"-"` // ignore flag: "-"
}

var u = &User{
    Name:"testing",
}

// Insert data with default driver.
if _, err := qsql.InsertStruct(mdb, u, "testing"); err != nil{
    // ... 
}
// ...

// Or Insert data with designated driver.
if _, err := qsql.InsertStruct(mdb, u, "testing", qsql.DRV_NAME_MYSQL); err != nil{
    // ... 
}
// ...
```

## Quick query way
``` text

// Way 1: query result to a struct.
type User struct{
    Id   int64 `db:"id"`
    Name string `db:"name"`
}

mdb := db.GetCache("master") 
// or mdb = <sql.Tx>
var u = *User{}
if err := qsql.QueryStruct(mdb, u, "SELECT id, name FROM a WHERE id = ?", id)
if err != nil{
    // sql.ErrNoRows has been replace by errors.ErrNoData
    if errors.ErrNoData.Equal(err) {
       // no data
    }
    // ...
}
// ..

// Way 2: query row to struct
mdb := db.GetCache("master") 
// or mdb = <sql.Tx>
var u = *User{}
if err := qsql.ScanStruct(qsql.QueryRow(mdb, "SELECT id, name FROM a WHERE id = ?", id), u); err != nil {
    // sql.ErrNoRows has been replace by errors.ErrNoData
    if errors.ErrNoData.Equal(err) {
       // no data
    }
    // ...
}

// Way 3: query result to structs
mdb := db.GetCache("master") 
// or mdb = <sql.Tx>
var u = []*User{}
if err := qsql.QueryStructs(mdb, &u, "SELECT id, name FROM a WHERE id = ?", id); err != nil {
    // ...
}
if len(u) == 0{
    // data not found
    // ...
}
// .. 

// Way 4: query rows to structs
mdb := db.GetCache("master") 
// or mdb = <sql.Tx>
rows, err := qsql.Query(mdb, "SELECT id, name FROM a WHERE id = ?", id)
if err != nil {
    // ...
}
defer qsql.Close(rows)
var u = []*User{}
if err := qsql.ScanStructs(rows, &u); err != nil{
    // ...
}
if len(u) == 0{
    // data not found
    // ...
}

```

## Query an element which is implemented sql.Scanner

```text
mdb := db.GetCache("master") 
// or mdb = <sql.Tx>
count := 0
if err := qsql.QueryElem(mdb, &count, "SELECT count(*) FROM a WHERE id = ?", id); err != nil{
    // sql.ErrNoRows has been replace by errors.ErrNoData
    if errors.ErrNoData.Equal(err) {
       // no data
    }
    // ...
}
```
## Extend the where in stmt
// Example for the first input:
mdb := db.GetCache("master") 
args:=[]int{1,2,3}
mdb.Query(fmt.Sprintf("select * from table_name where in (%s)", qsql.StmtWhereIn(0,len(args))), qsql.SliceToArgs(args)...)
// Or
mdb.Query(fmt.Sprintf("select * from table_name where in (%s)", qsql.StmtWhereIn(0,len(args), qsql.DRV_NAME_MYSQL), qsql.SliceToArgs(args)...)

// Example for the second input:
mdb.Query(fmt.Sprintf("select * from table_name where id=? in (%s)", qsql.StmtWhereIn(1,len(args)), id, qsql.SliceToArgs(args)...)

## Mass query.
```text
mdb := db.GetCache("master") 
qSql = &qsql.Page{
     CountSql:`SELECT count(1) FROM user_info WHERE create_time >= ? AND create_time <= ?`,
     DataSql:`SELECT mobile, balance FROM user_info WHERE create_time >= ? AND create_time <= ?`
}
count, titles, result, err := qSql.QueryPageArray(db, true, condition, 0, 10)
// ...
// Or
count, titles, result, err := qSql.QueryPageMap(db, true, condtion, 0, 10)
// ...
if err != nil {
// ...
}
```

## Make a MultiTx
``` text
multiTx := []*qsql.MultiTx{}
multiTx = append(multiTx, qsql.NewMultiTx(
    "UPDATE testing SET name = ? WHERE id = ?",
    id,
))
multiTx = append(multiTx, qsql.NewMultiTx(
    "UPDATE testing SET name = ? WHERE id = ?",
    id,
))

// do exec multi tx
mdb := db.GetCache("master") 
tx, err := mdb.Begin()
if err != nil{
    // ...
}
if err := qsql.ExecMutlTx(tx, multiTx); err != nil {
    qsql.Rollback(tx)
    // ...
}
if err := tx.Commit(); err != nil {
    qsql.Rollback(tx)
    // ...
}
```
