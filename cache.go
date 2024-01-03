package qsql

import (
	"os"
	"sync"
	"time"

	"github.com/go-ini/ini"
	"github.com/gwaylib/errors"
)

var (
	cacheLock    = sync.Mutex{}
	cacheIniPath string
	cache        = map[string]*DB{}
)

func setCacheIni(iniPath string) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cacheIniPath = iniPath
}

func regCache(key string, db *DB) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	_, ok := cache[key]
	if ok {
		panic("key is already exist: " + key)
	}
	cache[key] = db
}

func getCache(key string) (*DB, error) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	db, ok := cache[key]
	if !ok {
		if len(cacheIniPath) == 0 {
			return nil, errors.ErrNoData.As(key)
		}
		db, err := regCacheWithIni(cacheIniPath, key)
		if err != nil {
			return nil, errors.As(err)
		}
		return db, nil
	}
	return db, nil
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

// ini content example
//
// [main]
// driver: mysql
// dsn: username:passwd@tcp(127.0.0.1:3306)/main?timeout=30s&strict=true&loc=Local&parseTime=true&allowOldPasswords=1
// max_life_time:7200 # seconds
// max_idle_time:0 # seconds
// max_idle_conns:0 # num
// max_open_conns:0 # num
//
// [log]
// driver: mysql
// dsn: username:passwd@tcp(127.0.0.1:3306)/log?timeout=30s&strict=true&loc=Local&parseTime=true&allowOldPasswords=1
// max_life_time:7200 # seconds
//
func regCacheWithIni(iniFile, iniSection string) (*DB, error) {
	// create a new
	cfg, err := ini.Load(iniFile)
	if err != nil {
		return nil, errors.As(err, iniFile)
	}
	section, err := cfg.GetSection(iniSection)
	if err != nil {
		return nil, errors.As(err, iniSection)
	}
	drvName, err := section.GetKey("driver")
	if err != nil {
		return nil, errors.As(err, "not found 'driver'")
	}
	dsn, err := section.GetKey("dsn")
	if err != nil {
		return nil, errors.As(err, "not found 'dsn'")
	}

	// http://techblog.en.klab-blogs.com/archives/31093990.html
	lifeTime := int64(0)
	lifeTimeKey, err := section.GetKey("max_life_time")
	if err == nil {
		lifeTime, err = lifeTimeKey.Int64()
		if err != nil {
			return nil, errors.As(err, "error max_life_time value")
		}
	}
	idleTime := int64(0)
	idleTimeKey, err := section.GetKey("max_idle_time")
	if err == nil {
		idleTime, err = idleTimeKey.Int64()
		if err != nil {
			return nil, errors.As(err, "error max_idle_time value")
		}
	}

	idleConns := int(0)
	idleConnKey, err := section.GetKey("max_idle_conns")
	if err == nil {
		idleConns, err = idleConnKey.Int()
		if err != nil {
			return nil, errors.As(err, "error max_idle_conns value")
		}
	}

	openConns := int(0)
	openConnKey, err := section.GetKey("max_open_conns")
	if err == nil {
		openConns, err = openConnKey.Int()
		if err != nil {
			return nil, errors.As(err, "error max_open_conns value")
		}
	}

	_, ok := cache[iniSection]
	if ok {
		return nil, errors.New("key is already exist").As(iniSection)
	}

	db, err := Open(drvName.String(), os.ExpandEnv(dsn.String()))
	if err != nil {
		return nil, errors.As(err)
	}
	if lifeTime > 0 {
		db.SetConnMaxLifetime(time.Duration(lifeTime) * time.Second)
	}
	if idleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(idleTime) * time.Second)
	}
	if idleConns > 0 {
		db.SetMaxIdleConns(idleConns)
	}
	if openConns > 0 {
		db.SetMaxOpenConns(openConns)
	}

	cache[iniSection] = db
	return db, nil
}
