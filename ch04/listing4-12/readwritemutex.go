package main

import (
	"fmt"
	"sync"
	"time"
)

type ReadWriteMutex struct {
	readersCounter int
	readersLock    sync.Mutex
	globalLock     sync.Mutex
}

func (rw *ReadWriteMutex) ReadLock() {
	rw.readersLock.Lock()
	rw.readersCounter++
	if rw.readersCounter == 1 {
		rw.globalLock.Lock()
	}
	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
	rw.globalLock.Lock()
}

func (rw *ReadWriteMutex) ReadUnlock() {
	rw.readersLock.Lock()
	rw.readersCounter--
	if rw.readersCounter == 0 {
		rw.globalLock.Unlock()
	}
	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
	rw.globalLock.Unlock()
}

func main() {
	rwMutex := ReadWriteMutex{}
	for i := 0; i < 10; i++ {
		go func() {
			rwMutex.ReadLock()
			fmt.Println("Read started")
			time.Sleep(time.Second * 5)
			fmt.Println("Read finished")
			rwMutex.ReadUnlock()
		}()
	}
	time.Sleep(time.Second * 1)
	fmt.Println("Write started")
	rwMutex.WriteLock()
	fmt.Println("Write finished")
	rwMutex.WriteUnlock()
}
