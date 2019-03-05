package buffer

import (
	"testing"
	"time"
	"math/rand"
	"sync"
	"fmt"
)

func TestChan(t *testing.T) {
	c := NewChan()
	var wg sync.WaitGroup
	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			c.Put(1)
		}()
		go func() {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			c.Get()
			wg.Done()
		}()
	}
	wg.Wait()
}

func Test(t *testing.T) {
	c := make(chan int, 1)
	c <- 1
	close(c)
	v, ok := <-c
	fmt.Println(v, ok)
	v, ok = <-c
	fmt.Println(v, ok)
	v, ok = <-c
	fmt.Println(v, ok)
}