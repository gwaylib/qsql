package main

import (
	"fmt"

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

	// create table
	if _, err := mdb.Exec(
		`CREATE TABLE user (
		  "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  "username" VARCHAR(32) NOT NULL UNIQUE,
		  "passwd" VARCHAR(128) NOT NULL
		);`); err != nil {
		panic(err)
	}

	// std insert
	if _, err := mdb.Exec("INSERT INTO user(username,passwd)VALUES(?,?)", "t1", "t1"); err != nil {
		panic(err)
	}

	// reflect insert
	newUser := &TestingUser{UserName: "t2", Passwd: "t2"}
	if _, err := mdb.InsertStruct(newUser, "user"); err != nil {
		panic(err)
	}
	if newUser.ID == 0 {
		panic("expect newUser.ID > 0")
	}

	// std query
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

	// query where in
	whereIn := []string{"t1", "t2"}
	whereInCount := 0
	if err := mdb.QueryElem(&whereInCount,
		fmt.Sprintf("SELECT COUNT(*) FROM user WHERE username in (%s)", mdb.StmtWhereIn(0, len(whereIn))),
		qsql.StmtSliceArgs(whereIn)...,
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
	// TODO: more stmt optimization
	stmt, err := mdb.Prepare("SELECT COUNT(*) FROM user WHERE username=?")
	count := 0
	if err := stmt.QueryRow("t3").Scan(&count); err != nil {
		panic(err)
	}
	if count != 1 {
		panic(errors.New("need count==1").As(count))
	}
}
