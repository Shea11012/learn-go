package gee_cache

import (
	"unsafe"
)

// ByteView 缓存值，只读
// 使用 []byte 可以存储任意数据
type ByteView struct {
	b []byte
}

// Len 返回缓存值大小
func (b ByteView) Len() int {
	return len(b.b)
}

// ByteSlice 返回一个拷贝，防止其内部值被篡改
func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.b)
}

func (b ByteView) String() string {
	return byte2str(b.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte,len(b))
	copy(c,b)
	return c
}

func byte2str(b []byte) string {
	x := (*[2]uintptr)(unsafe.Pointer(&b))
	s := [2]uintptr{x[0],x[1]}
	return *(*string)(unsafe.Pointer(&s))
}

func str2byte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	b := [3]uintptr{x[0],x[1],x[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}
