# go_infrastructure

**用途:工具包，基于 golang 实现**

## auth 

- 基于 golang micro service 服务鉴权

```go
package main
import (
	"github.com/g-airport/go-infra/auth"
	"github.com/micro/go-micro"
)
var FuncRights = map[string]auth.Auth {
	// 控制日志是否输出
	"ImplSrvName.FuncName":{1,true,true},
}
    servive := micro.NewService()
    service.Init(
	// ...
	micro.WrapHandler(auth.ServerAuthWrapper(FuncRights)),
	)
```
- [Run Micro Service]

```go
//Prepare Interface
type ImplSrv struct {
}

func NewServer() *ImplSrv {
	return &ImplSrv{}
}

//Run go-micro
ns := NewServer()
server := service.Server()
server.Handle(server.NewHandler(ns))
```


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
consul (service discover) , mq , gRPC ...

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

## sync 

- [GlobalTimer]

- [Once]

## buffer

```go
import (
	"github.com/g-airport/go-infra/buffer"

)

    c := NewChan()
    c.Put(1)
    c.Get()
```

## Tool

| Usage 💡 | Link 🔑
| --- | --- |
|Golang embed Proxy   |[Golang embed Proxy](https://github.com/g-airport/go-infra/blob/master/proxy/readme.md) |
|Kafka Command     |[Kafka Command](https://github.com/g-airport/go-infra/blob/master/mq/readme.md) |
|Retry Func     |[Retry Func](https://github.com/g-airport/go-infra/blob/master/retry/retry.go) |
|Float Math     |[Float Math](https://github.com/g-airport/go-infra/blob/master/math/math.go) |
|Match Func     |[Match Func](https://github.com/g-airport/go-infra/blob/master/match/match.go) |
|User Context     |[User Context](https://github.com/g-airport/go-infra/blob/master/context/context.go) |
|AES crypt     |[AES crypt](https://github.com/g-airport/go-infra/blob/master/crypt/aes.go) |



