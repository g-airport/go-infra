package sync

import (
	"sync"
	"time"

	"github.com/micro/go-sync/lock"
	mConsul "github.com/micro/go-sync/lock/consul"

	"github.com/onlyLeoll/go-infra/log"
)

var global lock.Lock

type GlobalTimer struct {
	stop chan bool
	wg   sync.WaitGroup
}

func NewGlobalTimer() *GlobalTimer {
	return &GlobalTimer{
		stop: make(chan bool),
		wg:   sync.WaitGroup{},
	}
}

func (gt *GlobalTimer) Register(lock string, timeout time.Duration,
	f func() error) {
	go func() {
		gt.wg.Add(1)
		defer gt.wg.Done()

		for {
			select {
			case <-gt.stop:
				log.Info("%s loop stopped", lock)
				return
			case <-time.After(timeout):
				err := GlobalTransaction(lock, f)
				if err != nil {
					log.Err("%s failed: %v", lock, err)
				}
			}
		}
	}()
}

func (gt *GlobalTimer) Stop() {
	close(gt.stop)
	gt.wg.Wait()
	log.Info("global timer stopped")
}

func Init() {
	global = mConsul.NewLock()
}

func GlobalTransaction(lock string, f func() error) error {
	log.Info("global transaction: %v start", lock)
	err := global.Acquire(lock, nil)
	if err != nil {
		return err
	}

	defer func() {
		global.Release(lock)
		log.Info("global transaction: %v unlocked", lock)
	}()

	log.Info("global transaction: %v locked", lock)
	return f()
}

type Once struct {
	m map[string]bool
	l *sync.Mutex
}

func NewOnce() *Once {
	return &Once{
		m: make(map[string]bool),
		l: new(sync.Mutex),
	}
}

func (once *Once) Do(k string, f func() error) error {
	once.l.Lock()
	defer once.l.Unlock()

	_, ok := once.m[k]
	if ok {
		return nil
	}

	once.m[k] = true
	return f()
}
