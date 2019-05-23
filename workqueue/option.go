package workqueue

import "time"

type Option func(*WorkQueue)

func WithQueueName(name string) Option {
	return func(wq *WorkQueue) {
		wq.queueName = name
	}
}

func WithQueueSize(num int) Option {
	return func(wq *WorkQueue) {
		wq.queueSize = num
	}
}

func WithMaxWorker(member int) Option {
	return func(wq *WorkQueue) {
		wq.maxWorkers = member
	}
}

func WithMaxIdleWorker(num int) Option {
	return func(wq *WorkQueue) {
		wq.maxIdleWorkers = num
	}
}

func WithWatchPeriod(t time.Duration) Option {
	return func(wq *WorkQueue) {
		wq.watchPeriod = t
	}
}

func WithEnqueueTimeout(t time.Duration) Option {
	return func(wq *WorkQueue) {
		wq.enqueueTimeout = t
	}
}

func WithWorkerMaxWaiting(duration time.Duration) Option {
	return func(wq *WorkQueue) {
		wq.workerMaxWaiting = duration
	}
}

func WithHandler(h func(msg interface{})) Option {
	return func(wq *WorkQueue) {
		wq.handler = h
	}
}
