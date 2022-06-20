package geoip

import (
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

var once sync.Once
var _geoIp *geoip2.Reader

func Init(mmdbFile string) {
	if mmdbFile == "" {
		panic("not found mmdbFile")
	}
	once.Do(func() {
		var err error
		_geoIp, err = geoip2.Open(mmdbFile)
		if err != nil {
		}
	})
}

func GetRecordByIP(ipStr string) (*geoip2.City, error) {
	ip := net.ParseIP(ipStr)
	record, err := _geoIp.City(ip)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func Stop() {
	_ = _geoIp.Close()
}
