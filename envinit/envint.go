package envinit

import (
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	iCfg "github.com/g-airport/go-infra/config"
	iEnv "github.com/g-airport/go-infra/env"
	iLog "github.com/g-airport/go-infra/log"
	iSync "github.com/g-airport/go-infra/sync"

	"github.com/google/gops/agent"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/consul"
	_ "github.com/micro/go-plugins/client/grpc"
	_ "github.com/micro/go-plugins/server/grpc"
	//todo _ "github.com/micro/go-plugins/broker/nats"
)

var Debug = false

var (
	repository = "unknown"
	branch     = "unknown"
	commit     = "unknown"
)

var path = "G/api"

func Init() {
	cf := config.NewConfig(
		config.WithSource(
			consul.NewSource(
				consul.WithPrefix(path))))

	InitWithConfig(cf)

	iLog.Infow("service build on", "repository", repository, "branch", branch,
		"commit", commit)
}

func InitWithConfig(cf config.Config) {
	iCfg.Init(cf)
	Debug = iCfg.Debug()
	iLogLevel := iCfg.LogLevel()
	iLogger, err := iLog.NewLogger(iEnv.LogPath, iLogLevel)
	iEnv.ErrExit(err)

	iLog.SetDefault(iLogger)
	iSync.Init()
	InitPProfService()
}

func InitPProfService() {
	if Debug {
		runtime.SetBlockProfileRate(1)
	}

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		iLog.Info("start pprof service failed")
		return
	}

	iLog.Info("start pprof service on:", ln.Addr())
	go func() {
		http.Serve(ln, nil)
	}()

	if err := agent.Listen(agent.Options{}); err != nil {
		iLog.Info("start gops agent failed", err)
	}
}
