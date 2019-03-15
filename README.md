# go_infrastructure

**用途:工具包，基于 golang 实现**

## env 

- 目录配置

```golang
var (
	Dir        string
	RunDir     string
	LogDir     string
	LogPath    string
	ConfigPath string
	ConfigDir  string
	Pid        int
	Hostname   string
)
```

## config 

- init some infrastructure eg: some middleware like 
consul (service discover) , mq , gRpc ...

## log

- rotate log

```golang
import (
    "github.com/g-airport/go-infra/log"
    "github.com/g-airport/go-infra/env"
)
    
func init() {
    var Access log.Logger 
    path, _ := filepath.Abs(env.LogDir + "/access.log")
    A, err = glog.NewLogger(path, "debug")
    env.ErrExit(err)
    // 通过 var Access log.Logger 设定
    // 访问日志，调试日志，错误日志，... 
    // what ever you like 😁
}
//example
func log() {
    //hanler err
    var err error
    A.Info("", err)
}
```

## validate

- 检测 中文和可见字符

## errors

- 错误处理

```golang
type Error struct {
	Code     int         `json:"code"`
	Status   int         `json:"Status"`
	Detail   string      `json:"detail"`
	Internal string      `json:"internal,omitempty"`
	Content  interface{} `json:"content,omitempty"`
}
```

