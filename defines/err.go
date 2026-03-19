package defines

import (
	"errors"
)

var (
	ErrExist       = errors.New("target allready exists") // 目标已存在
	ErrNotExist    = errors.New("target do not exist")    // 目标不存在
	ErrUnSupported = errors.New("request is unsupported") // 要求不支持
	ErrInvalid     = errors.New("target invalid")         // 非法
)
