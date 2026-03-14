package gosdk

import (
	"time"
)

// 默认500ms周期调用 time.Now() 获取时间戳的时钟
//
// 用于非精确获取时间以及需求量大的场景
type CycleClock struct {
	Now        time.Time
	interval   time.Duration
	signalChan chan byte
}

// 周期时钟
func NewCycleClock() *CycleClock {
	return &CycleClock{
		interval:   500 * time.Millisecond,
		signalChan: make(chan byte),
	}
}

// 运行
func (c *CycleClock) Run() {
	c.Now = time.Now()
	go c.tick()
}

// 新周期
func (c *CycleClock) NewInterval(t time.Duration) {
	c.interval = t
	c.signalChan <- 0x01
}

// 关闭
func (c *CycleClock) Close() {
	c.signalChan <- 0x00
}

func (c *CycleClock) tick() {
	t := time.NewTicker(c.interval)
	for {
		select {
		case <-t.C:
			c.Now = time.Now()
		case b := <-c.signalChan:
			if b == 0x00 {
				return
			} else {
				t = time.NewTicker(c.interval)
			}
		}
	}
}
