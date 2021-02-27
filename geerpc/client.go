package geerpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"geerpc/message"
	"io"
	"log"
	"net"
	"sync"
)

type Call struct {
	Seq uint64
	ServiceMethod string // service.method
	Args interface{}	// 函数参数
	Reply interface{}	// 函数返回
	Error error
	Done chan *Call
}

func (c *Call) done() {
	c.Done <- c
}

type Client struct {
	message message.Message
	opt *Option
	sending sync.Mutex
	header message.Header
	mu sync.Mutex
	seq uint64
	pending map[uint64]*Call
	closing bool
	shutdown bool
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing {
		return ErrShutdown
	}

	c.closing = true
	return c.message.Close()
}

func (c *Client) IsAvailable() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return !c.shutdown && !c.closing
}

var _ io.Closer = (*Client)(nil)

var ErrShutdown = errors.New("connection is shutdown")

func (c *Client) registerCall(call *Call) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing || c.shutdown {
		return 0,ErrShutdown
	}

	call.Seq = c.seq
	c.pending[call.Seq] = call
	c.seq++
	return call.Seq,nil
}

func (c *Client) removeCall(seq uint64) *Call {
	c.mu.Lock()
	defer c.mu.Unlock()
	call := c.pending[seq]
	delete(c.pending,seq)
	return call
}

func (c *Client) terminateCalls(err error) {
	c.sending.Lock()
	defer c.sending.Unlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.shutdown = true
	for _,call := range c.pending {
		call.Error = err
		call.done()
	}
}

func (c *Client) receive() {
	var err error
	for err != nil {
		var h message.Header
		if err = c.message.ReadHeader(&h); err != nil {
			break
		}

		call := c.removeCall(h.Seq)
		switch {
		case call == nil:
			err = c.message.ReadBody(nil)
		case h.Error != "":
			call.Error = fmt.Errorf(h.Error)
			err = c.message.ReadBody(nil)
			call.done()
		default:
			err = c.message.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body" + err.Error())
			}
			call.done()
		}
	}
	c.terminateCalls(err)
}

func NewClient(conn net.Conn,opt *Option) (*Client,error) {
	f := message.NewMessageFuncMap[opt.MessageType]
	if f == nil {
		err := fmt.Errorf("invalid message type %s",opt.MessageType)
		log.Println("rpc client: message error:",err)
		return nil,err
	}

	if err := json.NewEncoder(conn).Encode(opt);err != nil {
		log.Println("rpc client: options error:",err)
		_ = conn.Close()
		return nil,err
	}
	return newClientMessage(f(conn),opt),nil
}

func newClientMessage(m message.Message,opt *Option) *Client {
	client := &Client{
		seq: 1,
		message: m,
		opt:opt,
		pending: make(map[uint64]*Call),
	}
	go client.receive()
	return client
}

func parseOptions(opts ...*Option) (*Option,error) {
	if len(opts) == 0|| opts[0] == nil {
		return DefaultOption,nil
	}

	if len(opts) != 1 {
		return nil,errors.New("number of options is more than 1")
	}

	opt := opts[0]
	opt.MagicNumber = DefaultOption.MagicNumber
	if opt.MessageType == "" {
		opt.MessageType = DefaultOption.MessageType
	}

	return opt,nil
}

func Dial(network,address string,opts ...*Option) (*Client,error) {
	opt,err := parseOptions(opts...)
	if err != nil {
		return nil,err
	}

	conn,err := net.Dial(network,address)
	if err != nil {
		return nil,err
	}
	client ,err := NewClient(conn,opt)
	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()

	return client,err
}