package main

import (
	"fmt"
	"sync"
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

func main() {
	semaphore := NewSemaphore(0)
	for i := 0; i < 5000; i++ {
		go doWork(semaphore)
		fmt.Println("Waiting for child goroutine")
		semaphore.Acquire()
		fmt.Println("Child goroutine finished")
	}
}

func doWork(semaphore *Semaphore) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	semaphore.Release()
}
