package wim

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Start(opt Option) error {
	r := mux.NewRouter()
	initPrometheus(r, opt.PrometheusOption)
	initBadger(r, opt.BadgerOption)
	return http.ListenAndServe(opt.Addr, r)
}
