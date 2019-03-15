package wim

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/dgraph-io/badger"
	"strings"
	"fmt"
)

func initBadger(r *mux.Router, opt *BadgerOption) {
	badgerRouter := r.PathPrefix("/badger").Subrouter()
	badgerRouter.HandleFunc("/{key}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		key, ok := vars["key"]
		if !ok {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("empty key"))
			return
		}
		if opt.DB == nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("nil db"))
			return
		}
		var sb strings.Builder
		count := 0
		err := opt.DB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				keyStr := string(item.Key())
				if strings.Contains(keyStr, key) {
					err := item.Value(func(val []byte) error {
						_, err := fmt.Fprintf(&sb, "key: %s, value: %s\n", keyStr, string(val))
						return err
					})
					if err != nil {
						return err
					}
					count++
				}
			}
			return nil
		})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(fmt.Sprintf("count: %d\n", count)))
		writer.Write([]byte(sb.String()))
	}).Methods("GET")
}
