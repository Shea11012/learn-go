package server

import (
	"chatroom/logic"
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 新用户加入到聊天室
	nickname := r.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal：", nickname)
		wsjson.Write(r.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度：4-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal")
		return
	}

	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("昵称已存在：", nickname)
		wsjson.Write(r.Context(), conn, logic.NewErrorMessage("该昵称已存在"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists")
		return
	}

	user := logic.NewUser(conn, nickname, r.RemoteAddr)
	// 给用户发送消息
	go user.SendMessage(r.Context())

	// 给新用户发送欢迎消息
	user.Welcome()

	// 告知所有用户新用户到来
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)

	// 将该用户加入到广播列表
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, "join chat")

	// 接受用户消息
	err = user.ReceiveMessage(r.Context())

	// 用户离开
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
	log.Println("user:", nickname, "leaves chat")

	// 根据不同的错误，返回不同的状态码
	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client err:", err)
		conn.Close(websocket.StatusInternalError, "read from client err")
	}

}
