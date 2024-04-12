package main

import (
	"fmt"
	"sync"
)

type WaitGrp struct {
	groupSize int
	cond      *sync.Cond
}

func NewWaitGrp() *WaitGrp {
	return &WaitGrp{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (wg *WaitGrp) Add(delta int) {
	wg.cond.L.Lock()
	wg.groupSize += delta
	wg.cond.L.Unlock()
}

func (wg *WaitGrp) Wait() {
	wg.cond.L.Lock()
	for wg.groupSize > 0 {
		wg.cond.Wait()
	}
	wg.cond.L.Unlock()
}

func (wg *WaitGrp) Done() {
	wg.cond.L.Lock()
	wg.groupSize--
	if wg.groupSize == 0 {
		wg.cond.Broadcast()
	}
	wg.cond.L.Unlock()
}

func doWork(id int, wg *WaitGrp) {
	fmt.Println(id, "Done working ")
	wg.Done()
}

func main() {
	wg := NewWaitGrp()
	for i := 1; i <= 4; i++ {
		wg.Add(2)
		go doWork(i, wg)
		go doWork(i, wg)
	}
	wg.Wait()
	fmt.Println("All complete")
}
