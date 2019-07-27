package graceful

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func TestWaitGroupGracefulShutdown(t *testing.T) {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		AddOne()
		go func() {
			defer Done()
			time.Sleep(time.Second * time.Duration(rand.Int63n(5)))
		}()
	}
	Shutdown(ctx)
}

func TestWaitGroupGracefulShutdownWithTimeOut(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	for i := 0; i < 10; i++ {
		AddOne()
		go func() {
			defer Done()
			time.Sleep(time.Second * time.Duration(rand.Int63n(5)))
		}()
	}
	Shutdown(ctx)
}

