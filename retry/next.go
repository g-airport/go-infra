package retry

import (
	"math"
	"sync"
	"time"
)

type Next interface {
	Next() bool
}

// sync retry
func Run(r Next, f func() error, onError func(err error), onPanic func(x interface{})) bool {
	var done bool
	var wg sync.WaitGroup
	for r.Next() {
		wg.Add(1)
		go func() {
			defer func() {
				if x := recover(); x != nil {
					onPanic(x)
				}
				wg.Done()
			}()
			if err := f(); err != nil {
				onError(err)
				return
			}
			done = true
		}()
		wg.Wait()
		if done {
			return true
		}
	}
	return false
}

// async retry
func RunAsync(r Next, f func() error, onError func(err error), onPanic func(x interface{})) chan bool {
	c := make(chan bool, 1)
	go func() {
		c <- Run(r, f, onError, onPanic)
	}()
	return c
}

// back off mechanism
// 1.after fixed time
// 2.pow 10 millisecond
type WaitFunc func() time.Duration

func NewFixedWaitFunc(d time.Duration) WaitFunc {
	return func() time.Duration {
		return d
	}
}

// 0 10 100 1000 ...
func NewBackOffWaitFunc() WaitFunc {
	var attempts = -1
	return func() time.Duration {
		attempts++
		if attempts == 0 {
			return 0
		}
		return time.Duration(math.Pow(10, float64(attempts))) * time.Millisecond
	}
}

type Counter struct {
	Count int64    // retry times
	Wait  WaitFunc //mechanism
	count int64    // count = Count ?
}

func (c *Counter) Next() bool {
	if c.count == c.Count {
		return false
	}
	if c.count > 0 {
		time.Sleep(c.Wait())
	}
	c.count++
	return true
}
