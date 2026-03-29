package gosdk

import (
	"runtime"
	"sync/atomic"
)

// 基于原子量的自旋锁
type SpinLock struct {
	value int32
}

func (lock *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32(&lock.value, 0, 1) {
		runtime.Gosched()
	}
}

func (lock *SpinLock) Unlock() {
	atomic.StoreInt32(&lock.value, 0)
}
