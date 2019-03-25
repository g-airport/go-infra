package http

import (
	"net/http"
	"time"
)

type Option struct {
	ServiceName      string
	Handler          http.Handler
	Version          string
	MetaData         map[string]string
	RegisterTTL      time.Duration
	RegisterInterval time.Duration
}
