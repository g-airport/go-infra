package geoip

import (
	"testing"
)

func Test_IP(t *testing.T) {
	Init("~/GeoLite2-Country.mmdb")
	r, err := GetRecordByIP("114.114.114.114")
	t.Log(r.Country.IsoCode, err)
	r, err = GetRecordByIP("8.8.8.8")
	t.Log(r.Country.IsoCode, err)
}
