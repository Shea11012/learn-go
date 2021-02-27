package geerpc

import (
	"encoding/json"
	"fmt"
	"geerpc/message"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int
	MessageType message.Type
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	MessageType: message.GobType,
}

type Server struct {

}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (s *Server) Accept(lis net.Listener) {
	for {
		conn,err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:",err)
			return
		}

		go s.ServeConn(conn)
	}
}

func (s *Server) ServeConn(conn io.ReadWriteCloser)  {
	defer func() {
		_ = conn.Close()
	}()

	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: option error",err)
		return
	}

	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x",opt.MagicNumber)
		return
	}

	f := message.NewMessageFuncMap[opt.MessageType]
	if f == nil {
		log.Printf("rpc server: invalid message Type %s",opt.MessageType)
		return
	}

	s.serveMessage(f(conn))
}

// invalidRequest 无效的请求
var invalidRequest = struct {}{}

func (s *Server) serveMessage(m message.Message) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	for {
		req,err := s.readRequest(m)
		if err != nil {
			if req == nil {
				break
			}

			req.h.Error = err.Error()
			s.sendResponse(m,req.h,invalidRequest,sending)
			continue
		}

		wg.Add(1)
		go s.handleRequest(m,req,sending,wg)
	}

	wg.Wait()
	_ = m.Close()
}

func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

type request struct {
	h *message.Header
	argv,replyv reflect.Value
}

func (s *Server) readRequestHeader(m message.Message) (*message.Header, error) {
	var h message.Header
	if err := m.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:",err)
		}

		return nil,err
	}

	return &h,nil
}

func (s *Server) readRequest(m message.Message) (*request, error) {
	h,err := s.readRequestHeader(m)
	if err != nil {
		return nil,err
	}

	req := &request{h: h}
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = m.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:",err)
	}

	return req,nil
}

func (s *Server) sendResponse(m message.Message, h *message.Header, body interface{},sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := m.Write(h,body);err != nil {
		log.Println("rpc server: write response error:",err)
	}
}

func (s *Server) handleRequest(m message.Message,req *request,sending *sync.Mutex,wg *sync.WaitGroup)  {
	defer wg.Done()
	log.Println(req.h,req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d",req.h.Seq))
	s.sendResponse(m,req.h,req.replyv.Interface(),sending)
}