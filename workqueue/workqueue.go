package workqueue

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/g-airport/go-infra/log"
)

var (
	ErrWorkQueueStopped = errors.New("work-queue has stopped")
	ErrEnqueueTimeout   = errors.New("enqueue timeout")

	defaultEnqueueTimeout   = time.Minute
	defaultWatchPeriod      = time.Second * 10 // print log
	defaultWorkerMaxWaiting = time.Minute * 10
	//defaultLoopTimeout      = time.Millisecond * 100 time tick
)

// base CSP
// step.1 feed data into queue
// step.2 deal data with handler func

type WorkQueue struct {
	queueName      string
	curWorkers     int
	maxWorkers     int
	maxIdleWorkers int
	workerID       int
	queueSize      int
	running        bool

	workerChan        chan *Worker
	workerStoppedChan chan struct{}
	stop              chan struct{}
	jobChan           chan interface{} // data length

	stopDispatchChan chan struct{}

	wg      *sync.WaitGroup
	handler func(msg interface{})

	watch *time.Ticker

	watchPeriod      time.Duration
	enqueueTimeout   time.Duration
	workerMaxWaiting time.Duration
}

func NewQueue(name string, timeout time.Duration, queueSize, workSize int,
	handler func(msg interface{})) *WorkQueue {
	return newQueue(
		WithQueueName(name),
		WithHandler(handler),
		WithMaxWorker(workSize),
		WithQueueSize(queueSize),
		WithEnqueueTimeout(timeout),
		WithMaxIdleWorker(maxIdleWorkerSize(workSize)),
	)
}

func NewQueueWithWatch(name string, timeout, watch time.Duration,
	queueSize, workSize int, handler func(msg interface{})) *WorkQueue {
	return newQueue(
		WithQueueName(name),
		WithHandler(handler),
		WithMaxWorker(workSize),
		WithQueueSize(queueSize),
		WithEnqueueTimeout(timeout),
		WithMaxIdleWorker(maxIdleWorkerSize(workSize)),
		WithWatchPeriod(watch),
	)
}

func NewQueueWithOptions(opts ...Option) *WorkQueue {
	return newQueue(opts...)
}

func maxIdleWorkerSize(i int) int {
	return int(math.Max(float64(i/2), 1.0))
}

func newQueue(opts ...Option) *WorkQueue {
	q := &WorkQueue{
		wg:               &sync.WaitGroup{},
		stop:             make(chan struct{}),
		workerMaxWaiting: defaultWorkerMaxWaiting,
		enqueueTimeout:   defaultEnqueueTimeout,
		watchPeriod:      defaultWatchPeriod,
		stopDispatchChan: make(chan struct{}),
	}

	for _, o := range opts {
		o(q)
	}

	q.jobChan = make(chan interface{}, q.queueSize)
	q.workerChan = make(chan *Worker, q.maxIdleWorkers)
	q.workerStoppedChan = make(chan struct{}, q.maxWorkers)
	q.init()
	return q
}

func (wq *WorkQueue) EnqueueWithTimeout(msg interface{}, timeout time.Duration) error {
	return wq.enqueue(msg, timeout)
}

func (wq *WorkQueue) Enqueue(msg interface{}) error {
	return wq.enqueue(msg, wq.enqueueTimeout)
}

func (wq *WorkQueue) enqueue(msg interface{}, t time.Duration) error {
	if !wq.running {
		return ErrWorkQueueStopped
	}

	timeout := time.NewTimer(t)
	defer timeout.Stop()

	select {
	case wq.jobChan <- msg:
		return nil
	case <-timeout.C:
		return ErrEnqueueTimeout
	}
}

func (wq *WorkQueue) Stop() {
	if !wq.running {
		return
	}

	wq.running = false
	wq.stopDispatchChan <- struct{}{}

	start := time.Now()
	log.Info("workQueue %v begin stop left items: %v", wq.queueName, len(wq.jobChan))
	defer func() {
		log.Info("workQueue %v stop cost %v", wq.queueName, time.Since(start))
	}()
	wq.wg.Wait()
	wq.watch.Stop()
}

func (wq *WorkQueue) init() {
	wq.running = true
	wq.watch = time.NewTicker(wq.watchPeriod)
	for i := 0; i < wq.maxIdleWorkers; i++ {
		wq.addWorker()
	}

	go wq.dispatch()
}

func (wq *WorkQueue) addWorker() *Worker {
	worker := NewWorker(wq.workerID, wq)
	wq.workerID++
	wq.curWorkers++
	return worker
}

func (wq *WorkQueue) dispatch() {
	for wq.running || len(wq.jobChan) > 0 || len(wq.workerStoppedChan) > 0 {
		select {
		case msg := <-wq.jobChan:
			worker := wq.getIdleWorker()
			worker.do(msg)
		case <-wq.workerStoppedChan:
			wq.curWorkers--
		case <-wq.watch.C:
			size := len(wq.jobChan)
			if size == 0 {
				continue
			}

			log.Info("workerQueue %v monitor running left items: %v currently workers:%v",
				wq.queueName, size, wq.curWorkers)

		case <-wq.stopDispatchChan:
		}
	}
	close(wq.stop)
}

func (wq *WorkQueue) getIdleWorker() *Worker {
	select {
	case w := <-wq.workerChan:
		return w
	default:
		if wq.curWorkers < wq.maxWorkers {
			wq.addWorker()
		}
	}
	w := <-wq.workerChan

	return w
}
