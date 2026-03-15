package gosdk

import "unsafe"

// 调用unsafe包实现
func Unsafe_Bs2Str(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	return unsafe.String(&bs[0], len(bs))
}

// 调用unsafe包实现
func Unsafe_Str2Bs(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
