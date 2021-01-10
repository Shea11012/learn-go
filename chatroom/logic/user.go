package logic

import (
	"chatroom/global"
	"context"
	"errors"
	"regexp"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type User struct {
	UID            int       `json:"uid"`
	NickName       string    `json:"nick_name"`
	EnterAt        time.Time `json:"enter_at"`
	Addr           string    `json:"addr"`
	messageChannel chan *Message
	conn           *websocket.Conn
}

func NewUser(conn *websocket.Conn, nickname string, add string) *User {
	return &User{
		NickName:       nickname,
		Addr:           add,
		conn:           conn,
		EnterAt:        time.Now(),
		messageChannel: make(chan *Message, global.MessageQueueLen),
	}
}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.messageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

func (u *User) ReceiveMessage(ctx context.Context) error {
	var receiveMsg map[string]string
	var err error
	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}
			return err
		}

		sendMsg := NewMessage(u, receiveMsg["content"], "")
		reg := regexp.MustCompile(`@[^\s@]{2,30}`)
		sendMsg.Ats = reg.FindAllString(sendMsg.Content, -1)
		Broadcaster.Broadcast(sendMsg)
	}
}

func (u *User) Welcome() {
	u.messageChannel <- NewWelcomeMessage(u)
}

func (u *User) CloseMessageChannel() {
	close(u.messageChannel)
}

func (u *User) WriteMessage(msg *Message) {
	u.messageChannel <- msg
}
