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
				server.File("main.go").Write(`package main

import (
    "` + name + `/internal/config"
	"path/filepath"
)

func main(){
    config.OnTest() // 开启Test

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
			component := internal.Dir("component") // 具体组件
			{
				component.File("database.go").Write(`package component

import (
	"github.com/wednis/gosdk"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	db, err := gosdk.NewMysqlGorm("root", "root", "test", nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}
`)
				component.File("logger.go").Write(`package component

import (
	"` + name + `/internal/config"

	"github.com/wednis/gosdk"
	"go.uber.org/zap"
)

func NewLogger(cfg *config.Config) *zap.Logger {
	return gosdk.NewZapLogger(cfg.IsTest())
}
`)
				component.File("validator.go").Write(`package component

import "github.com/go-playground/validator/v10"

func NewValidator() *validator.Validate {
	return validator.New()
}
`)
			}
			config := internal.Dir("config") // 读取配置文件数据
			{
				config.File("config.go").Write(`package config

import (
	"github.com/wednis/gosdk"
)

type ExternalConfig struct {
    // to be completed
}

type InternalConfig struct {
	// to be completed
}

type GlobalData struct {
	// to be completed
}

type TestConfig struct {
	// to be completed
}

type Config struct {
	External ExternalConfig
	Internal InternalConfig
	Global   GlobalData
	test     *TestConfig
}

func (cfg *Config) OnTest() {
	cfg.test = &TestConfig{
		// to be completed
	}
}

func (cfg *Config) IsTest() bool {
	return cfg.test != nil
}

func (cfg *Config) TestConfig() *TestConfig {
	return cfg.test
}

func InitConfig(path string) (*Config, error) {
	cfg := &Config{
		External: ExternalConfig{
			// to be completed
		},
		Internal: InternalConfig{
			// to be completed
		},
		Global: GlobalData{
			// to be completed
		},
		test: nil,
	}
	return cfg, gosdk.BindConfig(path, &cfg.External)
}
`)
			}
			internal.Dir("dto")                // 数据传输层
			handler := internal.Dir("handler") // HTTP请求处理器
			{
				handler.Dir("request")              // 请求数据相关
				response := handler.Dir("response") // 响应数据相关
				{
					response.File("response").Write(`package response
import "` + name + `/pkg/errcode"

type Response struct {
	Code errcode.Code ` + "`" + `json:"code"` + "`" + `
	Data any          ` + "`" + `json:"data"` + "`" + `
}
`)
				}
			}
			internal.Dir("service")          // 业务逻辑层
			internal.Dir("repository")       // 数据访问层
			internal.Dir("model")            // 数据库数据模型
			routes := internal.Dir("routes") // 路由注册
			{
				routes.File("routes.go").Write(`package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(
	engine *gin.Engine,
	// to be completed
) {
	// to be completed
}
`)
			}
			internal.Dir("middleware") // 中间件
		}
		pkg := folder.Dir("pkg")
		{
			pkg.Dir("utils")              // 通用工具函数
			errcode := pkg.Dir("errcode") // 错误码
			{
				errcode.File("errcode.go").Write(`package errcode

type Code int

const (
	Success Code = 0
)
`)
			}
		}
		script := folder.Dir("script") // 存放脚本
		{
			script.File("errcode.py") // 错误码生成脚本
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

// 在可执行文件所在的目录下生成预设的目录及文件
//   - config 存放配置文件
//   - log 存放日志
func GenExecFileDir() error {
	// TODO
	return nil
}
