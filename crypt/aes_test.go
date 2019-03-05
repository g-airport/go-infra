package crypt

import (
	"encoding/base64"
	"testing"
)

func TestAes(t *testing.T) {
	// AES-128
	// key length ：16, 24, 32 bytes ->  AES-128, AES-192, AES-256
	key := []byte("qwerasdfzxcvtyui")
	result, err := AesEncrypt([]byte("gggggg"), key)
	if err != nil {
		panic(err)
	}
	t.Log(base64.StdEncoding.EncodeToString(result))
	origData, err := AesDecrypt(result, key)
	if err != nil {
		panic(err)
	}
	t.Log(string(origData))
}

func TestData(t *testing.T) {
	//
	key := "q"
	t.Log(len(key))
	data := "Mozilla/5.0"
	ret, err := AesEncrypt([]byte(data), []byte(key))
	if err != nil {
		t.Log(ret)
	}
	_ = ret
}
