package gosdk

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/wednis/gosdk/defines"
	"gopkg.in/yaml.v3"
)

// map绑定到结构体
func bindFromMap(m map[string]any, result any) error {
	config := &mapstructure.DecoderConfig{
		Result: result,
		MatchName: func(mapKey, fieldName string) bool {
			return strings.EqualFold(mapKey, fieldName)
		}, // 忽略大小写区别
		ErrorUnused:      false, // 允许 map 中有未匹配字段
		WeaklyTypedInput: true,  // 允许类型转换（如字符串 "123" → int）
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(m)
}

func bindYAMLConfig(data []byte, cfgptr any) error {
	var m map[string]any
	if err := yaml.Unmarshal(data, &m); err != nil {
		return err
	}
	return bindFromMap(m, cfgptr)
}

func bindJSONConfig(data []byte, cfgptr any) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	return bindFromMap(m, cfgptr)
}

// 读取配置文件
//   - path 配置文件路径
//   - cfgptr 配置结构体指针（需要字段全部大写）
func BindConfig(path string, cfgptr any) error {
	if !IsFileExists(path) {
		return defines.ErrNotExist
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	switch filepath.Ext(path) {
	case ".json":
		return bindJSONConfig(data, cfgptr)
	case ".yaml":
		return bindYAMLConfig(data, cfgptr)
	default:
		return defines.ErrUnSupported
	}
}
