package registry

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RegistryConfig struct {
	HeartBeatInterval time.Duration // 心跳检测间隔（针对服务）
}

// 服务注册中心
type Registry struct {
	Port uint   // 开放端口
	Key  string // 用于确保Discovery和Registry之间建立连接

	engine   *gin.Engine
	lock     sync.RWMutex
	tokenMap map[string][]byte // 已注册服务Token映射表
}

func New(port uint, key string) *Registry {
	r := &Registry{
		Port: port,
		Key:  key,
	}
	r.init()
	return r
}

// 运行
func (r *Registry) Run() error {

}

func (r *Registry) init() {
	r.engine = gin.New()
	engine := r.engine
	engine.Use(r.auth)
	engine.GET("/heartbeat", heartbeatAPI(r))
	engine.POST("/register", registerAPI(r))
	engine.POST("/update", updateAPI(r))
	engine.GET("/services/:name", servicesAPI(r))
}
