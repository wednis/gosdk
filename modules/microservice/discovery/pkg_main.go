package discovery

import (
	"sync"
	"time"
)

type DiscoveryConfig struct {
	HeartBeatInterval time.Duration // 心跳检测间隔（针对Registry）
}

// 服务发现器
type Discovery struct {
	RegPort uint
	Key     string

	lock  sync.RWMutex        // 保护cache
	cache map[string][]string // 本地缓存
}

// 运行
func (d *Discovery) Run() {

}

// 获取服务地址
func (d *Discovery) Get(service string) string {
}
