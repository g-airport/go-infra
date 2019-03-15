package kv

import (
	"time"

	"github.com/dgraph-io/badger"
)

var Badger *badger.DB

// Entrance
func InitBadger() {
	var err error
	Badger, err = getBadger()
	if err != nil {
		panic(err)
	}
}

func getBadger() (*badger.DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/badger"
	opts.ValueDir = "/tmp/badger"
	db, err := badger.Open(opts)

	if err != nil {
		return nil, err
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
	}()

	return db, nil
}

// defer close conn
func Close() {
	_ = Badger.Close()
}
