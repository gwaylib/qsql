package qsql

import (
	"sync"
)

var (
	cacheLock = sync.Mutex{}
	cache     = map[string]*DB{}
)

func regCache(key string, db *DB) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	_, ok := cache[key]
	if ok {
		panic("key is already exist: " + key)
	}
	cache[key] = db
}

func getCache(key string) (*DB, bool) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	db, ok := cache[key]
	return db, ok
}

func rmCache(src *DB) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	for key, db := range cache {
		if src == db {
			delete(cache, key)
			return
		}
	}
}

func closeCache() {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	for key, db := range cache {
		Close(db)
		delete(cache, key)
	}
}
