package trace

import (
	"google.golang.org/grpc"
)

func GRPCServerWrapper()  {
	s := grpc.NewServer()
	// just use context ?

}