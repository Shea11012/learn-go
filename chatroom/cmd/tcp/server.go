package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	enteringChannel = make(chan *User)
	leavingChannel  = make(chan *User)
	messageChannel  = make(chan *Message, 8)
)

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	go broabcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		go handleConn(conn)
	}
}

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func (u *User) String() string {
	return "Addr: " + u.Addr + "UID: " + strconv.Itoa(u.ID) + "enterAt: " + u.EnterAt.Format("2006-01-02 15:04:02")
}

type Message struct {
	OwnerID int
	Content string
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	go sendMessage(conn, user.MessageChannel)

	user.MessageChannel <- "welcome, " + user.String()
	messageChannel <- &Message{
		OwnerID: user.ID,
		Content: user.String(),
	}
	enteringChannel <- user

	var userActive = make(chan struct{})
	go func() {
		d := 5 * time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChannel <- &Message{
			OwnerID: user.ID,
			Content: input.Text(),
		}
		userActive <- struct{}{}
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误:", err)
	}

	leavingChannel <- user
	messageChannel <- &Message{
		OwnerID: user.ID,
		Content: "user:" + strconv.Itoa(user.ID) + " has left",
	}
}

func sendMessage(conn net.Conn, channel chan string) {
	for msg := range channel {
		_, _ = fmt.Fprintln(conn, msg)
	}
}

var (
	globalID int
	idLocker sync.Mutex
)

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()
	globalID++
	return globalID
}

func broabcaster() {
	users := make(map[*User]struct{})
	for {
		select {
		case user := <-enteringChannel:
			users[user] = struct{}{}
		case user := <-leavingChannel:
			delete(users, user)
			close(user.MessageChannel)
		case msg := <-messageChannel:
			for user := range users {
				if user.ID != msg.OwnerID {
					user.MessageChannel <- msg.Content
				}
			}
		}
	}
}
