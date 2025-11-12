qsql is a supplement to the go sql package

# Refere to:
```
database/sql
https://github.com/jmoiron/sqlx
```

# Example:
More example see the [example](./example) directory.

## Directing use
``` text
package main

import (
	"github.com/gwaylib/conf"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/qsql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
    mdb, err := qsql.Open(qsql.DRV_NAME_MYSQL, dsn)
    if err != nil{
        panic(err)
    }
    arr := make([]string, 3)
    if err := mdb.QueryElems(&arr, "SELECT id, created_at, updated_at WHERE id=?", 1); err != nil{
        panic(err)
    }
}
```

## Using ini cache
the configuration file path like : './etc/db.cfg'

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

Make a package for cache.go with ini
``` text
package db

import (
	"github.com/gwaylib/conf"
	"github.com/gwaylib/qsql"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
    qsql.RegCacheWithIni(conf.RootDir() + "/etc/db.cfg")

    // Register cache without ini
    // db, err := qsql.Open(qsql.DRV_NAME_MYSQL, dsn)
    // if err != nil{
    //     panic(err)
    // }
    // qsql.RegCache("main", db)
}

func GetCache(section string) *qsql.DB {
	return qsql.GetCache(section)
}

func HasCache(section string) (*qsql.DB, error) {
	return qsql.HasCache(section)
}

func CloseCache() {
	qsql.CloseCache()
}
```

Using the cache package
``` text
package main

import (
	"github.com/gwaylib/conf"
	"github.com/gwaylib/errors"
	_ "github.com/go-sql-driver/mysql"

    "model/db"
)

func main() {
    mdb := db.GetCache("main")
    arr := make([]string, 3)
    if err := mdb.QueryElems(&arr, "SELECT id, created_at, updated_at WHERE id=?", 1); err != nil{
        panic(err)
    }
}
```

## Standard sql
*qsql.DB has implements *sql.DB, so you can call qsql.DB like *sql.DB
``` text
mdb := db.GetCache("main") 

row := mdb.QueryRow("SELECT * ...")
// ...

rows, err := mdb.Query("SELECT * ...")
// ...

result, err := mdb.Exec("UPDATE ...")
// ...
```

## Insert struct(s) into table
the struct tag format like `db:"field"`, reference to: http://github.com/jmoiron/sqlx
``` text
type User struct{
    Id     int64  `db:"id,auto_increment"` // flag "autoincrement", "auto_increment" are supported .
    Name   string `db:"name"`
    Ignore string `db:"-"` // ignore flag: "-"
}

func main() {
    mdb := db.GetCache("main") 

    var u = &User{
        Name:"testing",
    }

    // Insert data with driver.
    if _, err := mdb.InsertStruct(u, "testing"); err != nil {
        // ... 
    }
    // ...

    // Insert structs
    // TODO: qsql.InsertStructs ?
    var us := []User{}
    tx := mdb.Begin()
    txFn := func() error{
        for _, u := range us {
            if err := tx.InsertStruct(u, "testing"); err != nil{
                return errors.As(err)
            }
        }
        return nil
    }
    if err := qsql.Commit(tx, txFn); err != nil {
        // ...
    }
}

```

## Quick sql way
``` text
package main

import (
    gErrors "github.com/gwaylib/errors"
)

// Way 1: query result to a struct.
type User struct{
    Id   int64 `db:"id"`
    Name string `db:"name"`
}

func main() {
    mdb := db.GetCache("main") 
    var u = *User{}
    if err := mdb.QueryStruct(u, "SELECT id, name FROM a WHERE id = ?", id); err != nil {
        // sql.ErrNoRows has been replace by gErrors.ErrNoData
        if gErrors.ErrNoData.Equal(err) {
           // no data
        }
        // ...
    }
    // ..

    count := 0
    if err := mdb.QueryElem(&count, "SELECT count(*) FROM a WHERE id = ?", id); err != nil {
        // sql.ErrNoRows has been replace by errors.ErrNoData
        if errors.ErrNoData.Equal(err) {
           // no data
        }
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
    if err := mdb.Commit(tx, fn); err != nil {
        // ...
    }
}
```

## SqlBuilder
```text
func main() {
    mdb := qsql.GetCache("main") 

    id := 0
    inIds := []interface{}{1,2}

    bd := mdb.NewSqlBuilder() // qsql.NewSqlBuilder(mdb.DriverName())
    bd.Select("id", "created_at")
    bd.Add("FROM")
    bd.TabAdd("tmp")
    bd.Add("WHERE")
    bd.TabAdd("created_at BETWEEN ? AND ?", time.Now().AddDate(-1,0,0), time.Now())
    bd.IfTabAdd(len(inIds)>0, fmt.Sprintf("AND id IN (%s)", bd.In(inIds))) // create the sql params as '?', becareful, it has performance issues when the array is too long.
    titles, data, err := mdb.QueryPageArr(bd.String(), bd.Args()...) 
    if err != nil {
        panic(err)
    }

    updateBD := mdb.NewSqlBuilder()
    updateBD.Add("UPDATE tmp SET")
    updateBD.TabAdd("(updated_at=?,name=?)", time.Now())
    updateDB.Add("WHERE")
    updateDB.TabAdd("id=?", id)
    if _, err := mdb.Exec(updateDB.String(), updateDB.Args()...); err != nil {
        panic(err)
    }
}
```
