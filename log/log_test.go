package log

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/onlyLeoll/go-infra/env"
)

var Access *Logger

func init() {

	var err error
	path, _ := filepath.Abs(env.LogDir + "/access.log")
	fmt.Println(path)
	Access, err = NewLogger(path, "debug")
	if err != nil {
		panic("new log failed")
	}
}

func TestLog(t *testing.T) {
	Access.Info("access")
}
