package message

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobMessage struct {
	conn io.ReadWriteCloser
	buf *bufio.Writer
	dec *gob.Decoder
	enc *gob.Encoder
}

func (g *GobMessage) Close() error {
	return g.conn.Close()
}

func (g *GobMessage) ReadHeader(header *Header) error {
	return g.dec.Decode(header)
}

func (g *GobMessage) ReadBody(i interface{}) error {
	return g.dec.Decode(i)
}

func (g *GobMessage) Write(header *Header, body interface{}) (err error) {
	defer func() {
		_ = g.buf.Flush()
		if err != nil {
			_ = g.Close()
		}
	}()

	if err := g.enc.Encode(header);err != nil {
		log.Println("rpc codec:gob error encoding header:",err)
		return err
	}

	if err := g.enc.Encode(body);err != nil {
		log.Println("rpc codec: gob error encoding body:",err)
		return err
	}

	return nil
}

var _ Message = (*GobMessage)(nil)

func NewGobMessage(conn io.ReadWriteCloser) Message {
	buf := bufio.NewWriter(conn)
	return &GobMessage{
		conn: conn,
		buf: buf,
		dec: gob.NewDecoder(conn),
		enc: gob.NewEncoder(buf),
	}
}