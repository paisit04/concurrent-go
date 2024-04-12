package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	permits int
	cond    *sync.Cond
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		permits: n,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (rw *Semaphore) Acquire() {
	rw.cond.L.Lock()
	for rw.permits <= 0 {
		rw.cond.Wait()
	}
	rw.permits--
	rw.cond.L.Unlock()
}

func (rw *Semaphore) Release() {
	rw.cond.L.Lock()

	rw.permits++
	rw.cond.Signal()

	rw.cond.L.Unlock()
}

type Channel[M any] struct {
	capacitySema *Semaphore
	sizeSema     *Semaphore
	mutex        sync.Mutex
	buffer       *list.List
}

func NewChannel[M any](capacity int) *Channel[M] {
	return &Channel[M]{
		capacitySema: NewSemaphore(capacity),
		sizeSema:     NewSemaphore(0),
		buffer:       list.New(),
	}
}

func (c *Channel[M]) Send(message M) {
	c.capacitySema.Acquire()

	c.mutex.Lock()
	c.buffer.PushBack(message)
	c.mutex.Unlock()

	c.sizeSema.Release()
}

func (c *Channel[M]) Receive() M {
	c.capacitySema.Release()

	c.sizeSema.Acquire()

	c.mutex.Lock()
	v := c.buffer.Remove(c.buffer.Front()).(M)
	c.mutex.Unlock()

	return v
}

func receiver(messages *Channel[int], wGroup *sync.WaitGroup) {
	msg := 0
	for msg != -1 {
		time.Sleep(1 * time.Second)
		msg = messages.Receive()
		fmt.Println("Received:", msg)
	}
	wGroup.Done()
}

func main() {
	channel := NewChannel[int](10)
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go receiver(channel, &wGroup)
	for i := 1; i <= 6; i++ {
		fmt.Println("Sending: ", i)
		channel.Send(i)
	}
	channel.Send(-1)
	wGroup.Wait()
}
