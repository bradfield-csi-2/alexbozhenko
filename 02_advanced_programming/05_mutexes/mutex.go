package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	locked   int32 = 1
	unlocked int32 = 0
)

type myMutex struct {
	locked int32 // 1 -locked, 0 -unlocked
}

func (m *myMutex) Lock() {
	// if locked, then busywait
	for !atomic.CompareAndSwapInt32(&(m.locked), unlocked, locked) {
		// and try to lock atomically
		// func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
		// if the value in memory matches old, write new into addr
		//sleep
		time.Sleep(10 * time.Millisecond)
	}
}

func (m *myMutex) Unlock() {
	// if mutex is already unlocked, it is a runtime error

	if !atomic.CompareAndSwapInt32(&(m.locked), locked, unlocked) {
		panic("trying to unlock mutex which is already unlocked")
	}

}

const (
	numGoroutines = 1000
	numIncrements = 1000
)

var globalLock myMutex

type counter struct {
	count int
}

func safeIncrement(c *counter) {
	// globalLock.Lock()
	// defer globalLock.Unlock()

	// c.count += 1

	// cond.L.lock
	//

	c.count += 1
	//cond.L.unlock
	//cond.signal
}

const (
	mutexLocked = 1 << iota // mutex is locked

	mutexWoken

	mutexStarving

	mutexWaiterShift = iota
)

// changed to global mutex usage
func main() {
	fmt.Println(mutexLocked, mutexWoken, mutexStarving, mutexWaiterShift)
	c := &counter{
		count: 0,
	}

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < numIncrements; j++ {
				safeIncrement(c)
			}
		}()
	}

	wg.Wait()
	//fmt.Println(c.count)
}
