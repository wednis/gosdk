package gosdk

import "os"

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

// 0755权限
func NewRWXDir(path string) error {
	return os.Mkdir(path, 0755)
}

// 0755权限
func NewRWXFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
}
