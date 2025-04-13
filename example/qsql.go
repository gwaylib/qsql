package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/qsql"
	_ "github.com/mattn/go-sqlite3"
)

type TestingUser struct {
	ID       int64  `db:"id,auto_increment"` // auto_increment or autoincrement
	UserName string `db:"username"`
	Passwd   string `db:"passwd"`
}

func main() {
	mdb, _ := qsql.Open("sqlite3", ":memory:")
	defer qsql.Close(mdb)

	// create table by manualy
	if _, err := mdb.Exec(
		`CREATE TABLE user (
		  "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  "updated_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  "username" VARCHAR(32) NOT NULL UNIQUE,
		  "passwd" VARCHAR(128) NOT NULL
		);`); err != nil {
		panic(err)
	}

	// std sql insert one user
	if _, err := mdb.Exec("INSERT INTO user(username,passwd)VALUES(?,?)", "t1", "t1"); err != nil {
		panic(err)
	}

	// reflect insert one user
	newUser := &TestingUser{UserName: "t2", Passwd: "t2"}
	if _, err := mdb.InsertStruct(newUser, "user"); err != nil {
		panic(err)
	}
	if newUser.ID == 0 {
		panic("expect newUser.ID > 0")
	}

	// std sql query
	var id int64
	var username, passwd string
	if err := mdb.QueryRow("SELECT id, username, passwd FROM user WHERE username=?", "t1").Scan(&id, &username, &passwd); err != nil {
		panic(err)
	}
	if username != "t1" && passwd != "t1" {
		panic(username + "," + passwd)
	}
	if id == 0 {
		panic(id)
	}

	// reflect query
	// query struct data
	expectUser := &TestingUser{}
	if err := mdb.QueryStruct(expectUser, "SELECT * FROM user WHERE username=?", "t1"); err != nil {
		panic(err)
	}
	if expectUser.UserName != "t1" && expectUser.Passwd != "t1" {
		panic("data not match")
	}
	users := []*TestingUser{}
	if err := mdb.QueryStructs(&users, "SELECT * FROM user LIMIT 2"); err != nil {
		panic(err)
	}
	if len(users) != 2 {
		panic("expect len==2")
	}

	// query elememt data
	pwd := ""
	if err := mdb.QueryElem(&pwd, "SELECT passwd FROM user WHERE username=?", "t1"); err != nil {
		panic(err)
	}
	if pwd != "t1" {
		panic(pwd)
	}
	ids := []int64{}
	if err := mdb.QueryElems(&ids, "SELECT id FROM user LIMIT 2"); err != nil {
		panic(err)
	}
	if len(ids) != 2 {
		panic("expect len==2")
	}
	fmt.Printf("ids:%+v\n", ids)

	// query where if condition
	ifBD := mdb.NewSqlBuilder()
	ifBD.Select("id,created_at").
		Add("FROM user WHERE id=?", "t1").
		AddIf(rand.Int()%2 == 0, "OR (created_at BETWEN ? AND ?)", time.Now().Add(-1e9), time.Now())
	if _, _, err := mdb.QueryPageArr(ifBD.String(), ifBD.Args()...); err != nil {
		panic(err)
	}

	// query where in
	whereIn := []string{"t1", "t2"}
	whereInCount := 0
	sqlbd := mdb.NewSqlBuilder()
	sqlbd.Select("COUNT(*)").
		Add("FROM").
		AddTab("user").
		Add("WHERE").
		AddTab("username in (" + sqlbd.AddStmtIn(whereIn) + ")")

	if err := mdb.QueryElem(&whereInCount,
		sqlbd.String(),
		sqlbd.Args()...,
	); err != nil {
		panic(err)
	}
	if whereInCount != 2 {
		panic("expect count of whereIn is 2")
	}

	// query data in string
	// table type
	titles, data, err := mdb.QueryPageArr("SELECT * FROM user LIMIT 10")
	if err != nil {
		panic(err)
	}
	fmt.Printf("PageArr title:%+v\n", titles)
	fmt.Printf("PageArr data: %+v\n", data)
	// map type
	titles, mData, err := mdb.QueryPageMap("SELECT * FROM user LIMIT 10")
	if err != nil {
		panic(err)
	}
	fmt.Printf("PageMap title:%+v\n", titles)
	fmt.Printf("PageMap data: %+v\n", mData)

	// executer for tx
	tx, err := mdb.Begin()
	if err != nil {
		panic(err)
	}
	if err := mdb.Commit(tx, func() error {
		txUsers := []TestingUser{
			{UserName: "t3", Passwd: "t3"},
			{UserName: "t4", Passwd: "t4"},
		}
		for _, u := range txUsers {
			if _, err := qsql.InsertStruct(mdb.DriverName(), tx, &u, "user"); err != nil {
				return errors.As(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	// excute for stmt
	stmt, err := mdb.Prepare(mdb.NewSqlBuilder().
		Select("COUNT(*)").Add("FROM user WHERE username=?").String(),
	)
	count := 0
	if err := stmt.QueryRow("t3").Scan(&count); err != nil {
		panic(err)
	}
	if count != 1 {
		panic(errors.New("need count==1").As(count))
	}

	// excute update
	updateBD := mdb.NewSqlBuilder().
		Add("UPDATE user SET passwd=? WHERE id=?", "t3", "t3")
	if _, err := mdb.Exec(updateBD.String(), updateBD.Args()...); err != nil {
		panic(errors.As(err, updateBD.Sql))
	}
}
