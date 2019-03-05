package buffer

import "sync"

type Chan struct {
	c      chan interface{}
	m      sync.Mutex
	buffer []interface{}
}

func NewChan() *Chan {
	return &Chan{
		c: make(chan interface{}, 1),
	}
}

func (c *Chan) Put(v interface{}) {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.buffer) == 0 {
		select {
		case c.c <- v:
			return
		default:
		}
	}
	c.buffer = append(c.buffer, v)
}

func (c *Chan) load() {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.buffer) > 0 {
		select {
		case c.c <- c.buffer[0]:
			c.buffer = c.buffer[1:]
			c.buffer[0] = nil // gc
		default:
		}
	}
}

func (c *Chan) Get() interface{} {
	v := <-c.c
	c.load()
	return v
}