package gosdk

import "reflect"

// 基于反射的操作

// 是不是结构体指针
func IsStructPointer(v any) bool {
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct
}

// 是不是方法
func IsFunction(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Func
}

// 是不是切片
func IsSlice(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Slice
}

// 是不是映射表
func IsMap(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Map
}

// 是不是通道
func IsChan(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Chan
}

// 是不是数组
func IsArray(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Array
}
