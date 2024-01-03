qsql is a supplement to the go sql package

# Refere to:
```
database/sql
https://github.com/jmoiron/sqlx
```

# Example:
More example see the [example](./example) directory.

## Using a manual cache
``` text
package db

import (
	"github.com/gwaylib/conf"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/qsql"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
    db, err := qsql.Open(qsql.DRV_NAME_MYSQL, dsn)
    if err != nil{
        panic(err)
    }
    qsql.RegCache("main", db)
}
```


## Using etc cache
Assume that the configuration file path is: './etc/db.cfg'

The etc file content
```
[main]
driver: mysql
dsn: username:passwd@tcp(127.0.0.1:3306)/main?timeout=30s&strict=true&loc=Local&parseTime=true&allowOldPasswords=1
max_life_time:7200 # seconds
max_idle_time:0 # seconds
max_idle_conns:0 # num
max_open_conns:0 # num

[log]
driver: mysql
dsn: username:passwd@tcp(127.0.0.1:3306)/log?timeout=30s&strict=true&loc=Local&parseTime=true&allowOldPasswords=1
max_life_time:7200
```

Make a package for connection cache with ini
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
    if err := qsql.RegCacheWithIni(dbFile, "main"); err != nil{
       panic(err)
    }
    if err := qsql.RegCacheWithIni(dbFile, "log"); err != nil{
       panic(err)
    }
}
```

Call a cache
``` text
mdb := qsql.GetCache("main")
```

## Call standar sql
``` text
mdb := qsql.GetCache("main") 

row := mdb.QueryRow("SELECT * ...")
// ...

rows, err := mdb.Query("SELECT * ...")
// ...

result, err := mdb.Exec("UPDATE ...")
// ...
```

## Insert a struct to db(using reflect)
``` text
type User struct{
    Id     int64  `db:"id,auto_increment"` // flag "autoincrement", "auto_increment" are supported .
    Name   string `db:"name"`
    Ignore string `db:"-"` // ignore flag: "-"
}

func main() {
    mdb := qsql.GetCache("main") 

    var u = &User{
        Name:"testing",
    }

    // Insert data with driver.
    if _, err := mdb.InsertStruct(u, "testing"); err != nil{
        // ... 
    }
    // ...
}

```

## Quick query way
``` text

// Way 1: query result to a struct.
type User struct{
    Id   int64 `db:"id"`
    Name string `db:"name"`
}

func main() {
    mdb := qsql.GetCache("main") 
    var u = *User{}
    if err := mdb.QueryStruct(u, "SELECT id, name FROM a WHERE id = ?", id)
    if err != nil{
        // sql.ErrNoRows has been replace by errors.ErrNoData
        if errors.ErrNoData.Equal(err) {
           // no data
        }
        // ...
    }
    // ..
    
    // Way 2: query row to struct
    mdb := qsql.GetCache("main") 
    var u = *User{}
    if err := mdb.ScanStruct(mdb.QueryRow("SELECT id, name FROM a WHERE id = ?", id), u); err != nil {
        // sql.ErrNoRows has been replace by errors.ErrNoData
        if errors.ErrNoData.Equal(err) {
           // no data
        }
        // ...
    }
    
    // Way 3: query result to structs
    mdb := qsql.GetCache("main") 
    var u = []*User{}
    if err := mdb.QueryStructs(&u, "SELECT id, name FROM a WHERE id = ?", id); err != nil {
        // ...
    }
    if len(u) == 0{
        // data not found
        // ...
    }
    // .. 
    
    // Way 4: query rows to structs
    mdb := qsql.GetCache("main") 
    rows, err := mdb.Query("SELECT id, name FROM a WHERE id = ?", id)
    if err != nil {
        // ...
    }
    defer qsql.Close(rows)
    var u = []*User{}
    if err := mdb.ScanStructs(rows, &u); err != nil{
        // ...
    }
    if len(u) == 0{
        // data not found
        // ...
    }

}

```

## Query an element which is implemented sql.Scanner

```text
func main() {
    mdb := qsql.GetCache("main") 
    count := 0
    if err := mdb.QueryElem(&count, "SELECT count(*) FROM a WHERE id = ?", id); err != nil{
        // sql.ErrNoRows has been replace by errors.ErrNoData
        if errors.ErrNoData.Equal(err) {
           // no data
        }
        // ...
    }
}
```
## Extend the where in stmt
```text
// Example for the first input:
func main() {
    mdb := qsql.GetCache("main") 
    args:=[]int{1,2,3}
    mdb.Query(fmt.Sprintf("select * from table_name where in (%s)", mdb.StmtWhereIn(0, len(args))), qsql.StmtSliceArgs(args)...)
    // Or
    mdb.Query(fmt.Sprintf("select * from table_name where in (%s)", mdb.StmtWhereIn(0, len(args), qsql.DRV_NAME_MYSQL), qsql.StmtSliceArgs(args)...)
    
    // Example for the second input:
    mdb.Query(fmt.Sprintf("select * from table_name where id=? in (%s)", qsql.StmtWhereIn(1,len(args)), qsql.StmtSliceArgs(id, args)...)
}
```

## Mass query.
```text
func main() {
    mdb := qsql.GetCache("main") 
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
}
```

## Make a lazy tx commit
``` text
// commit the tx
func main() {
    mdb := qsql.GetCache("main") 
    tx, err := mdb.Begin()
    if err != nil{
        // ...
    }
    fn := func() error {
      if err := tx.Exec("UPDATE testing SET name = ? WHERE id = ?", id); err != nil{
        return err
      }
      return nil
    }
    if err := qsql.Commit(tx, fn); err != nil {
        // ...
    }
}
```
