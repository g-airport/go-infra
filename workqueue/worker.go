package workqueue

import (
	"time"

	"github.com/g-airport/go-infra/log"
)

type Worker struct {
	id      int
	job     chan interface{}
	queue   *WorkQueue
	running bool
}

func NewWorker(id int, wq *WorkQueue) *Worker {
	w := &Worker{
		id:      id,
		job:     make(chan interface{}),
		queue:   wq,
		running: true,
	}

	go w.loop()
	return w
}

func (w *Worker) stop() {
	w.running = false
	w.queue.workerStoppedChan <- struct{}{}
}

func (w *Worker) loop() {
	w.queue.wg.Add(1)
	defer w.queue.wg.Done()

	for w.running {
		w.idling()
		if !w.running {
			return
		}
		w.working()
	}
}

// worker recv job
func (w *Worker) do(i interface{}) {
	w.job <- i
}

// worker doing handle
func (w *Worker) working() {
	defer func() {
		if err := recover(); err != nil {
			log.Err("workQueue %v-%v panic by: %v stack: %v", w.queue.queueName, w.id,
				err, StackRecord())
		}
	}()

	select {
	case msg := <-w.job:
		w.queue.handler(msg)
	case <-w.queue.stop:
		w.stop()
	}
}

func (w *Worker) idling() {
	t := time.NewTimer(w.queue.workerMaxWaiting)
	defer t.Stop()
	select {
	case w.queue.workerChan <- w:
	case <-t.C:
		w.stop()
	case <-w.queue.stop:
		w.stop()
	}
}
