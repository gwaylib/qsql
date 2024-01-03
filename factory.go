package qsql

// Register a db to the connection pool by manully.
func RegCache(key string, db *DB) {
	regCache(key, db)
}

func RegCacheWithIni(iniPath string) {
	setCacheIni(iniPath)
}

// Get the db instance from the cache.
func GetCache(key string) *DB {
	db, err := getCache(key)
	if err != nil {
		panic(err)
	}
	return db
}

// Checking the cache does it have a db instance.
func HasCache(key string) (*DB, error) {
	return getCache(key)
}

// Close all instance in the cache.
func CloseCache() {
	closeCache()
}
