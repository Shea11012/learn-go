package server

import (
	"chatroom/logic"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

func RegisterHandle() {
	inferRootDir()
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
	http.HandleFunc("/user_list", userListHandleFunc)
}

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(rootDir + "/template/home.html")
	if err != nil {
		fmt.Fprintf(w, "模板解析错误")
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Fprint(w, "模板执行错误")
		return
	}
}

func userListHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	userList := logic.Broadcaster.GetUserList()
	b, err := json.Marshal(userList)
	if err != nil {
		fmt.Fprintf(w, `[]`)
	} else {
		fmt.Fprintf(w, string(b))
	}
}

var rootDir string

func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(d string) string
	infer = func(d string) string {
		if exists(d + "/template") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	rootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
