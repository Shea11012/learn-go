package gee_cache

import (
	"unsafe"
)

type ByteView struct {
	b []byte
}

func (b ByteView) Len() int {
	return len(b.b)
}

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
