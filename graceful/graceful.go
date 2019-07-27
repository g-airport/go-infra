package graceful

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)


var DefaultWg = NewWaitGroup()

func AddOne() {
	DefaultWg.AddOne()
}

func Done() {
	DefaultWg.Done()
}

func Shutdown(ctx context.Context) {
	DefaultWg.Shutdown(ctx)
}

func IsShutting() bool {
	return DefaultWg.IsShutting()
}

func Shutting() <-chan struct{} {
	return DefaultWg.Shutting
}

// Is shutting
const inShutting = 1

// WaitGroup add graceful shutdown
type WaitGroup struct {
	wg      *sync.WaitGroup
	once    *sync.Once
	allDone chan struct{}
	// Unfinished (count running goroutine)
	Unfinished int64
	// Shutting notify running goroutine
	Shutting chan struct{}
	// isShutting listen shutdown
	isShutting int32
}

func NewWaitGroup() *WaitGroup {
	return &WaitGroup{
		wg:       new(sync.WaitGroup),
		once:     new(sync.Once),
		allDone:  make(chan struct{}),
		Shutting: make(chan struct{}),
	}
}

func (w *WaitGroup) Add(delta int64) {
	var p int64
	for {
		p = atomic.LoadInt64(&w.Unfinished)
		if atomic.CompareAndSwapInt64(&w.Unfinished, p, p+delta) {
			w.wg.Add(int(delta))
			return
		}
	}
}

func (w *WaitGroup) AddOne() {
	w.Add(1)
}

func (w *WaitGroup) Done() {
	w.Add(-1)
}

func (w *WaitGroup) Shutdown(ctx context.Context) {
	w.once.Do(
		func() {
			go func() {
				close(w.Shutting)
				atomic.StoreInt32(&w.isShutting, inShutting)
				defer close(w.allDone)
				w.wg.Wait()
			}()
		})
	start := time.Now()
	ticker := time.NewTicker(time.Millisecond * 100)
	counterP := atomic.LoadInt64(&w.Unfinished)
	for {
		select {
		case <-w.allDone:
			log.Printf("[graceful]Wait: %s , %d goroutine done", time.Since(start), counterP)
			return
		case <-ticker.C:
			log.Printf("[graceful]Waitting: %s, left %d/%d un done, wait...", time.Since(start), atomic.LoadInt64(&w.Unfinished), counterP)
		case <-ctx.Done():
			counterA := atomic.LoadInt64(&w.Unfinished)
			switch ctxErr := ctx.Err(); ctxErr {
			case context.DeadlineExceeded:
				log.Printf("[graceful]ctx time out ,time %s, left %d/%d un done", time.Since(start), counterA, counterP)
			case context.Canceled:
				log.Printf("[graceful]ctx canceled, time %s, left %d/%d un done", time.Since(start), counterA, counterP)
			default:
				log.Printf("[graceful]unknown ctx.Err:%T %v, time %s, left %d/%d un done", ctxErr, ctxErr, time.Since(start), counterA, counterP)
			}
			return
		}
	}
}

// IsShutting
func (w *WaitGroup) IsShutting() bool {
	return atomic.LoadInt32(&w.isShutting) == inShutting
}
