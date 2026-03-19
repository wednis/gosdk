package gosdk

import (
	"errors"
	"path"
	"unsafe"
)

// 生成代码相关

// 待生成文件
type GenFile struct {
	Path string // 绝对路径
	Data []byte // 内容
}

func (f *GenFile) Write(data any) error {
	switch v := data.(type) {
	case string:
		f.Data = unsafe.Slice(unsafe.StringData(v), len(v))
	case []byte:
		f.Data = v
	default:
		return errors.New("unsupported type")
	}
	return nil
}

func (f *GenFile) Gen() error {
	file, err := NewRWXFile(f.Path)
	if err != nil {
		return err
	}
	if f.Data != nil {
		_, err = file.Write(f.Data)
	}
	return err
}

// 待生成目录
type GenDir struct {
	Path  string
	Files map[string]*GenFile
	Dirs  map[string]*GenDir
}

func (d *GenDir) File(name string) *GenFile {
	f := &GenFile{Path: path.Join(d.Path, name)}
	if d.Files == nil {
		d.Files = make(map[string]*GenFile)
	}
	// 如果存在直接覆盖
	d.Files[name] = f
	return f
}

func (d *GenDir) Dir(name string) *GenDir {
	dir := &GenDir{Path: path.Join(d.Path, name)}
	if d.Dirs == nil {
		d.Dirs = make(map[string]*GenDir)
	}
	// 如果存在直接覆盖
	d.Dirs[name] = dir
	return dir
}

// 递归生成目录结构
func (d *GenDir) Gen() error {
	var err error
	err = NewRWXDir(d.Path)
	if err != nil {
		return err
	}
	for _, f := range d.Files {
		if err = f.Gen(); err != nil {
			return err
		}
	}
	for _, dir := range d.Dirs {
		if err = dir.Gen(); err != nil {
			return err
		}
	}
	return nil
}

func vscodeGoDevWeb(root string, name string) *GenDir {
	folder := &GenDir{Path: path.Join(root, name)}
	// TODO
	{
		_vscode := folder.Dir(".vscode")
		_vscode.File("settings.json").Write(`{
    "go.goroot": "${env:HOME}/dev/env/go/1.26.0",
    "go.toolsEnvVars": {
        "GOMODCACHE": "${env:HOME}/dev/env/go/1.26.0/mod",
        "GOBIN": "${env:HOME}/dev/env/go/1.26.0/bin",
        "GOPROXY": "https://goproxy.cn,direct"
    },
    "terminal.integrated.env.linux": {
        "GO111MODULE": "on",
        "GOBIN": "${env:HOME}/dev/env/go/1.26.0/mod/bin",
        "GOMODCACHE": "${env:HOME}/dev/env/go/1.26.0/mod",
        "GOPROXY": "https://goproxy.cn,direct",
        "GOPATH": "${env:HOME}/dev/env/go/1.26.0/mod",
        "GOROOT": "${env:HOME}/dev/env/go/1.26.0",
        "PATH": "${env:HOME}/dev/env/go/1.26.0/bin:${env:HOME}/dev/env/go/1.26.0/mod/bin:${env:PATH}"
    },
    "liveServer.settings.root": "/web/"
}`)
		_vscode.File("tasks.json").Write(`{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "build-go",
            "type": "shell",
            "command": "go build -o ./build/server ./cmd/server/main.go",
            "group": {
                "kind": "build",
                "isDefault": true
            },
            "options": {
                "cwd": "${workspaceFolder}"
            },
            "problemMatcher": [
                "$go"
            ]
        }
    ]
}
`)
		cmd := folder.Dir("cmd")
		{
			server := cmd.Dir("server")
			{
				/*
									server.File("inject.go").Write(`package main

					import "github.com/wednis/gosdk"

					func inject(deps ...any) {
						err := gosdk.Inject(deps...).Invoke(
						    // to be completed
						).Err()
						if err != nil{
						    panic("failed to inject dependencies, error: " + err.Error())
						}
					}
					`)
				*/
				server.File("main.go").Write(`package main

import (
    "` + name + `/internal/config"
	"path/filepath"
)

func main(){
    config.OnDebug() // 开启Debug

	execdir, err := gosdk.GetExecDir()
	if err != nil {
		panic(err.Error())
	}
	cfg, err := config.InitConfig(filepath.Join(execdir, "config"))
	if err != nil {
		panic(err.Error())
	}

    // to be completed
	inject()
	// to be completed
}
`)
			}
		}
		folder.Dir("build")  // 存放编译结果
		folder.Dir("config") // 配置文件存放
		internal := folder.Dir("internal")
		{
			component := internal.Dir("component") // 组件（数据库 redis等）
			{
				component.File("database.go").Write(`package component

import "gorm.io/gorm"

type DB struct {
	*gorm.DB
}

func NewDatabase() (*DB, error) {
    // to be completed
	return nil, nil
}

`)
				component.File("logger.go").Write(`package component

import (
	"` + name + `/internal/config"

	"github.com/wednis/gosdk"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func NewLogger() *Logger {
	return &Logger{gosdk.NewZapLogger(config.IsDebug())}
}
`)
				component.File("validator.go").Write(`package component

import "github.com/go-playground/validator/v10"

type Validator struct {
    *validator.Validate
}

func NewValidator() *Validator {
    return &Validator{validator.New()}
}
`)
			}
			config := internal.Dir("config") // 读取 /config 下的配置文件数据
			{
				config.File("config.go").Write(`package config

import (
	"sync"

	"github.com/wednis/gosdk"
)

type Config struct {
    // to be completed

	debug  bool
	locker sync.RWMutex
}

// 开启DEBUG
func (cfg *Config) OnDebug() {
	cfg.locker.Lock()
	cfg.debug = true
	cfg.locker.Unlock()
}

// 关闭DEBUG
func (cfg *Config) OffDebug() {
	cfg.locker.Lock()
	cfg.debug = false
	cfg.locker.Unlock()
}

// 获取DEBUG状态
func (cfg *Config) IsDebug() bool {
	cfg.locker.RLock()
	defer cfg.locker.RUnlock()
	return debug
}

func NewConfig(path string) (*Config, error) {
	cfg := &Config{
	    // to be completed
	}
	return cfg, gosdk.BindConfig(path, cfg)
}
`)
				config.File("debug.go").Write(`package config

var debug bool

func OnDebug() {
	debug = true
}

func OffDebug() {
	debug = false
}

func IsDebug() bool {
	return debug
}
`)
			}
			internal.Dir("handler")    // HTTP请求处理器
			internal.Dir("service")    // 业务逻辑层
			internal.Dir("repository") // 数据访问层
			internal.Dir("model")      // 数据模型
			internal.Dir("middleware") // 中间件
		}
		pkg := folder.Dir("pkg")
		{
			pkg.Dir("utils") // 通用工具函数
		}
		folder.Dir("test") // 测试文件
		folder.Dir("web")  // 前端项目
	}
	folder.File(".gitignore").Write(`.vscode/
build/
config/
test/
web/
`)
	folder.File("README.md").Write("# " + name + `
`)
	folder.File("go.mod").Write("module " + name + `

go 1.26
`)
	return folder
}

// 生成Vscode环境的GO编写的WEB应用的目录结构
func GenVscodeGoDevWeb(root string, name string) error {
	if IsDirExists(root) {
		return vscodeGoDevWeb(root, name).Gen()
	}
	return errors.New("root not exists")
}
