package message

import (
	"io"
)

type Header struct {
	ServiceMethod string // "Service.Method"
	Seq uint64	// 请求序号
	Error string
}

type Message interface {
	io.Closer
	ReadHeader(*Header)	error
	ReadBody(interface{}) error
	Write(*Header,interface{}) error
}

type NewMessageFunc func(io.ReadWriteCloser) Message

type Type string

const (
	GobType Type = "application/gob"
	JsonType Type = "application/json"
)

var NewMessageFuncMap map[Type]NewMessageFunc

func init() {
	NewMessageFuncMap = make(map[Type]NewMessageFunc)
	NewMessageFuncMap[GobType] = NewGobMessage
}
