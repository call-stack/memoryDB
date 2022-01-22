package internal

import (
	"sync"
)

type Database struct {
	values map[string]string
	mu     sync.RWMutex
}

func NewDatabase() Database {
	db := Database{}
	db.values = make(map[string]string)

	return db
}

func (db *Database) GetValue(key string) string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var result string
	if val, ok := db.values[key]; ok {
		result = val
	} else {
		result = "value not present in db."
	}

	return result
}

func (db *Database) SetValue(key string, val string) string {
	db.mu.Lock()
	defer db.mu.Unlock()
	var result string
	if _, ok := db.values[key]; !ok {
		db.values[key] = val
		result = "OK"
	} else {
		result = "key already present"

	}
	return result
}

func (db *Database) DeleteValue(key string) string {
	db.mu.Lock()
	defer db.mu.Unlock()
	var result string
	if _, ok := db.values[key]; ok {
		delete(db.values, key)
		result = "OK"
	} else {
		result = "value not present in database"
	}

	return result
}
