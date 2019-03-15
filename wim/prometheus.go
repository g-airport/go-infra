package wim

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gorilla/mux"
)

func initPrometheus(r *mux.Router, opt *PrometheusOption) {
	prometheusRouter := r.PathPrefix("/badger").Subrouter()
	prometheus.MustRegister(opt.Collectors...)
	if opt.HttpHandler == nil {
		prometheusRouter.Handle("/metrics", promhttp.Handler()).Methods("GET")
	} else {
		prometheusRouter.Handle("/metrics", opt.HttpHandler).Methods("GET")
	}
}
