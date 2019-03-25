package http

import (
	"errors"
	"log"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-plugins/server/http"
)

// start and register a micro http service
func StartHttpSrv(opt Option) error {
	if opt.Handler == nil {
		return errors.New("nil handler")
	}
	if opt.RegisterTTL < 0 {
		return errors.New("negative ttl")
	}

	srvOpts := make([]server.Option, 0)
	if opt.ServiceName != "" {
		srvOpts = append(srvOpts, server.Name(opt.ServiceName))
	}
	if opt.MetaData != nil {
		srvOpts = append(srvOpts, server.Metadata(opt.MetaData))
	}
	if opt.Version != "" {
		srvOpts = append(srvOpts, server.Version(opt.Version))
	}
	if opt.RegisterInterval <= 0 {
		opt.RegisterInterval = time.Second * 30
	}
	srvOpts = append(srvOpts, server.RegisterTTL(opt.RegisterTTL))

	serverHttp := http.NewServer(srvOpts...)
	err := server.Handle(server.NewHandler(opt.Handler))
	if err != nil {
		log.Println("new http handler failed", err)
	}

	service := micro.NewService()
	// init cmd
	service.Init()
	// override
	service.Init(micro.Server(serverHttp), micro.RegisterInterval(opt.RegisterInterval))
	return service.Run()
}
