package qsql

import (
	slog "log"
)

type Log interface {
	Println(msg ...interface{})
}

var log = Log(slog.Default())

func SetLog(l Log) {
	log = l
}
