package gosdk

import "sync/atomic"

// 基于原子量的自旋锁
type SpinLock struct {
	value int32
}

func (lock *SpinLock) Lock() {
	// 原子交换，0换成1
	for !atomic.CompareAndSwapInt32(&lock.value, 0, 1) {
		// TODO 让出时间片或者递增
	}
}

func (lock *SpinLock) Unlock() {
	// 原子置零
	atomic.StoreInt32(&lock.value, 0)
}
