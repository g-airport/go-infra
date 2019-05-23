package workqueue

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/g-airport/go-infra/log"
)

func TestQueueStopWait(t *testing.T) {
	f := func(msg interface{}) {
		time.Sleep(time.Millisecond * 400)
	}
	sleep := time.Second
	msgSize := 100
	q := initQueue(f, time.Millisecond*50, sleep/time.Duration(msgSize), "queue", msgSize)
	time.Sleep(time.Second)
	q.Stop()
}

func initQueue(h func(interface{}), timeout, sleep time.Duration, name string,
	size int) *WorkQueue {
	log.Stdout()
	q := NewQueueWithOptions(
		WithQueueName(name),
		WithHandler(h),
		WithMaxWorker(10),
		WithQueueSize(64),
		WithEnqueueTimeout(timeout),
		WithWorkerMaxWaiting(time.Second),
	)

	go func() {
		i := 0
		for ; i < size; i++ {
			err := q.Enqueue(i)
			if err != nil {
				fmt.Println(i, "enqueue failed:", err)
			}
		}
	}()

	return q
}

func TestQueueRecover(t *testing.T) {
	f := func(msg interface{}) {
		time.Sleep(time.Millisecond * 100)
		m := msg.(string)
		fmt.Println(m)
	}

	sleep := time.Second
	msgSize := 10
	q := initQueue(f, time.Millisecond*50, sleep/time.Duration(msgSize), "queue", msgSize)
	time.Sleep(time.Second)
	q.Stop()
}

func TestQueueFullTimeout(t *testing.T) {
	f := func(msg interface{}) {
		time.Sleep(time.Second * 1)
	}

	sleep := time.Second * 5
	msgSize := 100
	q := initQueue(f, time.Millisecond*50, time.Millisecond, "queue", msgSize)
	time.Sleep(sleep)
	q.Stop()
}

func initQueueWithIdleWorkerSize(h func(interface{}),
	timeout, waitExpire time.Duration,
	name string,
	size int,
	workerSize, idleWorkerSize int) *WorkQueue {
	log.Stdout()
	q := NewQueueWithOptions(
		WithQueueName(name),
		WithEnqueueTimeout(timeout),
		WithQueueSize(64),
		WithMaxWorker(workerSize),
		WithMaxIdleWorker(idleWorkerSize),
		WithHandler(h),
		WithWorkerMaxWaiting(waitExpire),
		WithWatchPeriod(waitExpire*5))

	go func() {
		i := 0
		for ; i < size; i++ {
			err := q.Enqueue(i)
			var sleep time.Duration
			if i%200 == 0 && i > 0 {
				sleep = time.Duration(rand.Int63n(int64(time.Second * 30)))
			}
			time.Sleep(sleep)
			if err != nil {
				fmt.Println(i, "enqueue failed:", err)
			}
		}
		log.Info("finished send msg")
	}()

	return q
}

func TestQueueDynamicScale(t *testing.T) {
	f := func(msg interface{}) {
		s := time.Duration(time.Second * 5)
		time.Sleep(s)
	}

	msgSize := 1000
	q := initQueueWithIdleWorkerSize(f, time.Minute, time.Second,
		"queue", msgSize, 130, 3)
	time.Sleep(time.Minute * 2)
	q.Stop()
}

func TestWorkQueue_Stop(t *testing.T) {
	log.Stdout()
	count := 0
	q := NewQueueWithOptions(
		WithQueueName("test"),
		WithMaxWorker(1),
		WithQueueSize(10),
		WithEnqueueTimeout(time.Minute*10),
		WithMaxIdleWorker(0),
		WithWorkerMaxWaiting(time.Minute),
		WithHandler(func(msg interface{}) {
			time.Sleep(time.Second)
			count++
		}),
	)
	for i := 0; i < 10; i++ {
		_ = q.Enqueue(nil)
	}
	q.Stop()
	if count != 10 {
		t.Error("stop error")
	}
}
