package client

import (
	"communication_/pkg/model"
	"encoding/json"
	"fmt"
	"net"
)

var user model.User

func HandleConn(conn net.Conn) {
	ch := make(chan bool)
	go Output(conn, ch)

	Input(conn, ch)
}
func Input(conn net.Conn, ch chan bool) {
	for {
		for {
			HandleLoginInput(conn)
			if <-ch {
				break
			}
		}
		HandleStorageMessage(conn)
		var operation string
		flag := true
		for flag {
			fmt.Println("输入操作类型")
			fmt.Scan(&operation)
			switch operation {
			case "broadcast":
				HandleBroadCastInput(conn)
			case "logged_list":
				HandleLoggedListInput(conn)
			case "private_chat":
				HandlePrivateChatInput(conn)
			case "logout":
				HandleLogOutInput(conn)
				flag = false
			default:
				fmt.Println("没有该操作")
			}
		}
	}
}
func Output(conn net.Conn, ch chan bool) {
	b := make([]byte, 1024)
	m := make(map[string]interface{})
	for {
		n, _ := conn.Read(b)
		json.Unmarshal(b[:n], &m)
		Type := m["type"]
		data := m["data"]
		code := int(m["code"].(float64))
		msg := m["msg"]
		switch Type {
		case "login":
			if code == 200 {
				HandleLoginOutput(data.(map[string]interface{}), ch)
				ch <- true
			} else {
				fmt.Println(msg)
				ch <- false
			}
		case "logged_list":
			HandleLoggedListOutput(data.(map[string]interface{}))
		case "message_receive":
			HandleMessageReceiveOutput(data.(map[string]interface{}))

		default:
			fmt.Println("不能解析返回的操作类型")
		}
	}
}
func HandleLoginInput(conn net.Conn) {
	var username, password string
	fmt.Printf("输入用户名 密码：\n")
	fmt.Scan(&username, &password)
	data := make(map[string]string)
	data["username"] = username
	data["password"] = password
	m := make(map[string]interface{})
	m["type"] = "login"
	m["data"] = data
	b, _ := json.Marshal(m)
	conn.Write(b)
	fmt.Printf("登录加载中\n")
}
func HandleLoginOutput(data map[string]interface{}, ch chan bool) {
	user.Id = data["id"].(string)
	user.Username = data["username"].(string)
	user.Password = data["password"].(string)
}
func HandleLoggedListInput(conn net.Conn) {
	m := make(map[string]interface{})
	m["type"] = "logged_list"
	b, _ := json.Marshal(m)
	conn.Write(b)
}
func HandleLoggedListOutput(data map[string]interface{}) {
	fmt.Printf("用户列表: \n")
	for _, username := range data["logged_list"].([]string) {
		fmt.Printf("用户名: %v\n", username)
	}
}
func HandleBroadCastInput(conn net.Conn) {
	var body string
	fmt.Printf("输入广播消息：\n")
	fmt.Scan(&body)

	message := model.Message{Type: "broadcast", Sender: user.Username, Body: body, Receiver: "all"}
	b, _ := json.Marshal(message)
	data := make(map[string]interface{})
	json.Unmarshal(b, &data)
	r := make(map[string]interface{})
	r["type"] = "message_send"
	r["data"] = data
	b, _ = json.Marshal(r)
	conn.Write(b)
}
func HandlePrivateChatInput(conn net.Conn) {
	fmt.Printf("输入用户名 消息内容：\n")
	var receiver, body string
	fmt.Scan(&receiver, &body)

	message := model.Message{Type: "private", Sender: user.Username, Body: body, Receiver: receiver}
	b, _ := json.Marshal(message)
	data := make(map[string]interface{})
	json.Unmarshal(b, &data)
	r := make(map[string]interface{})
	r["type"] = "message_send"
	r["data"] = data
	b, _ = json.Marshal(r)
	conn.Write(b)
}

func HandleMessageReceiveOutput(data map[string]interface{}) {
	switch data["type"] {
	case "broadcast":
		fmt.Printf("发送方:%v 广播消息: %v\n", data["sender"], data["body"])

	case "private":
		fmt.Printf("发送方:%v 私聊消息: %v\n", data["sender"], data["body"])

	case "group":
		fmt.Printf("发送方:%v 群聊消息: %v\n", data["sender"], data["body"])

	default:
		fmt.Println("没有该消息类型")

	}
}
func HandleStorageMessage(conn net.Conn) {
	data := make(map[string]interface{})
	data["receiver"] = user.Username
	r := make(map[string]interface{})
	r["type"] = "storage_message"
	r["data"] = data
	b, _ := json.Marshal(r)
	conn.Write(b)
}
func HandleLogOutInput(conn net.Conn) {
	data := make(map[string]interface{})
	data["username"] = user.Username
	r := make(map[string]interface{})
	r["type"] = "logout"
	r["data"] = data
	b, _ := json.Marshal(r)
	conn.Write(b)
}
