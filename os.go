package gosdk

import (
	"encoding/json"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/wednis/gosdk/defines"
	"gopkg.in/yaml.v3"
)

func isPathExists(path string, m int) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	if m == 1 {
		return true
	} else if m == 2 {
		return !info.IsDir()
	} else {
		return info.IsDir()
	}
}

// 路径是否存在
func IsPathExists(path string) bool {
	return isPathExists(path, 1)
}

// 文件是否存在
func IsFileExists(path string) bool {
	return isPathExists(path, 2)
}

// 目录是否存在
func IsDirExists(path string) bool {
	return isPathExists(path, 3)
}

// 0755权限 本用户可读写执行其他用户读写
func NewRWXDir(path string) error {
	return os.Mkdir(path, 0755)
}

// 0755权限 本用户可读写执行其他用户读写
// 创建路径不存在的所有目录 目录已存在不报错
func NewRWXDirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

// 0755权限 本用户可读写执行其他用户读写
func NewRWXFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
}

// 新配置文件
// 0644权限 本用户可读写其他用户只读
func NewConfigFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}

// 获取编译后的可执行文件绝对路径（会溯源软链接）
func GetExecPath() (string, error) {
	p, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(p)
}

// 获取编译后的可执行文件所在目录（会溯源软链接）
func GetExecDir() (string, error) {
	p, err := GetExecPath()
	return filepath.Dir(p), err
}

// 生成json文件 权限0644
func NewJsonFile(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "    ") // 缩进
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// 生成yaml文件 权限0644
func NewYamlFile(path string, v any) error {
	file, err := NewConfigFile(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	return encoder.Encode(v)
}

// 获取以B为单位的文件大小
func GetFileSizeB(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, defines.ErrNotExist
		}
		return 0, err
	}
	return fileInfo.Size(), nil
}

// 获取以KB为单位的文件大小，保留两位小数
func GetFileSizeKB(path string) (float64, error) {
	fileBSize, err := GetFileSizeB(path)
	if err != nil {
		return 0, err
	}
	fileMBSize := float64(fileBSize) / 1024
	return math.Round(fileMBSize*100) / 100, nil
}

// 获取以MB为单位的文件大小，保留两位小数
func GetFileSizeMB(path string) (float64, error) {
	fileBSize, err := GetFileSizeB(path)
	if err != nil {
		return 0, err
	}
	fileMBSize := float64(fileBSize) / (1024 * 1024)
	return math.Round(fileMBSize*100) / 100, nil
}

// 获取以GB为单位的文件大小，保留两位小数
func GetFileSizeGB(path string) (float64, error) {
	fileBSize, err := GetFileSizeB(path)
	if err != nil {
		return 0, err
	}
	fileGBSize := float64(fileBSize) / (1024 * 1024 * 1024)
	return math.Round(fileGBSize*100) / 100, nil
}

// 获取以TB为单位的文件大小，保留两位小数
func GetFileSizeTB(path string) (float64, error) {
	fileBSize, err := GetFileSizeB(path)
	if err != nil {
		return 0, err
	}
	fileGBSize := float64(fileBSize) / (1024 * 1024 * 1024 * 1024)
	return math.Round(fileGBSize*100) / 100, nil
}

// 等待退出信号出现后执行
func WaitExitSignal(fn func()) {
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		<-sigs
		signal.Stop(sigs)
		fn()
	}()
}

// 获取操作系统类型
func GetOsKind() string {
	return runtime.GOOS
}

// 获取操作系统目录分隔符
func GetOsSep() string {
	return string(os.PathSeparator)
}

// 获取工作目录
func GetWorkDir() (string, error) {
	return os.Getwd()
}

// 设置工作目录
func SetWorkDir(dir string) error {
	return os.Chdir(dir)
}
