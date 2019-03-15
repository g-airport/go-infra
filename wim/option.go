package wim

import (
	"net/http"

	"github.com/dgraph-io/badger"
	"github.com/prometheus/client_golang/prometheus"
)

type Option struct {
	Addr             string
	PrometheusOption *PrometheusOption
	BadgerOption     *BadgerOption
}

type PrometheusOption struct {
	Collectors  []prometheus.Collector
	HttpHandler http.Handler
}

type BadgerOption struct {
	DB *badger.DB
}

// ...