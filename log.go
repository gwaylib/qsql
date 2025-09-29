package qsql

import (
	l "github.com/gwaylib/log"
)

type Log interface {
	Debug(msg ...interface{})
	Error(msg ...interface{})
}

var log = Log(l.NewWithCaller("qsql", 4))

func SetLog(l Log) {
	log = l
}
