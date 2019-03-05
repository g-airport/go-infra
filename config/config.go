package config

// base on micro

import (
	"log"
	"time"

	"github.com/micro/go-config"
	"github.com/micro/go-config/reader"
)

var Cfg config.Config

func Init(c config.Config) {
	Cfg = c
}

func Debug() bool {
	return Cfg.Get("debug").Bool(true)
}

func LogLevel() string {
	return Cfg.Get("log_level").String("debug")
}

func APIHost() string {
	return Cfg.Get("api_host").String("0.0.0.0:8080")
}

func ApiAllowedOrigins() []string {
	return Cfg.Get("api_allowed_origins").StringSlice([]string{"*"})
}

func RegisterTTL() (time.Duration, time.Duration) {
	ttl := Cfg.Get("register_ttl").Duration(30 * time.Second)
	interval := ttl / 2
	return ttl, interval
}

// key : value
func Get(key string) string {
	return Cfg.Get(key).String(key)
}

//Watch is wrapper of go-config's watch
func Watch(key string, fn func(watcher reader.Value), retryTime int) {
	if retryTime <= 0 {
		retryTime = 10
	}
	w, err := Cfg.Watch(key)
	if err != nil {
		log.Printf("watch config failed key:%v err:%v", key, err)
	}

	for retryTime > 0 {
		v, err := w.Next()
		if err != nil {
			log.Printf("get config change value failed err:%v", err)
			retryTime--
			time.Sleep(time.Second * 10)
			continue
		}
		fn(v)
	}
}
