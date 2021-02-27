package main

import (
	"encoding/json"
	"fmt"
	"geerpc"
	"geerpc/message"
	"log"
	"net"
	"time"
)

func main() {
	addr := make(chan string)
	go startServer(addr)
	conn,_ := net.Dial("tcp",<-addr)
	defer func() {
		_ = conn.Close()
		close(addr)
	}()

	time.Sleep(time.Second)
	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	m := message.NewGobMessage(conn)
	for i:=0;i<5;i++ {
		h := &message.Header{
			ServiceMethod: "Foo.Sum",
			Seq: uint64(i),
		}

		_ = m.Write(h,fmt.Sprintf("geerpc req %d",h.Seq))
		_ = m.ReadHeader(h)
		var reply string
		_ = m.ReadBody(&reply)
		log.Println("reply:",reply)
	}
}

func startServer(addr chan string) {
	l,err := net.Listen("tcp",":0")
	if err != nil {
		log.Fatal("network error:",err)
	}

	log.Println("start rpc server on",l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}


