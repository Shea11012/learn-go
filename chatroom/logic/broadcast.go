package logic

import (
	"chatroom/global"
)

var Broadcaster = &broadcaster{
	users:                 make(map[string]*User),
	enteringChannel:       make(chan *User),
	leavingChannel:        make(chan *User),
	messageChannel:        make(chan *Message, global.MessageQueueLen),
	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

type broadcaster struct {
	// 所有聊天室用户
	users map[string]*User

	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	// 判断昵称是否重复，是否可进入聊天室
	checkUserChannel      chan string
	checkUserCanInChannel chan bool

	requestUserChannel chan struct{}
	usersChannel       chan []*User
}

func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			b.users[user.NickName] = user
			b.sendUserList()
		case user := <-b.leavingChannel:
			delete(b.users, user.NickName)
			user.CloseMessageChannel()
			b.sendUserList()
		case msg := <-b.messageChannel:
			for i := range b.users {
				user := b.users[i]
				if user.UID == msg.User.UID {
					continue
				}
				user.WriteMessage(msg)
			}
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case <-b.requestUserChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList
		}
	}
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname // 因为该通道是无缓冲的所以可以避免使用锁

	return <-b.checkUserCanInChannel
}

func (b *broadcaster) sendUserList() {

}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	b.messageChannel <- msg
}

func (b *broadcaster) GetUserList() []*User {
	b.requestUserChannel <- struct{}{}
	return <-b.usersChannel
}
