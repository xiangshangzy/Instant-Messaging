package server

import (
	"communication_/pkg/model"
	"encoding/json"
	"fmt"
	"net"
)

var ConnLoggedMap = make(map[string]net.Conn)
var StorageMessageMap = make(map[string][]map[string]interface{})

func HandleConn(conn net.Conn) {
	var b = make([]byte, 1024)
	for {
		n, err := conn.Read(b)
		if err != nil {
			fmt.Println("客户端强制退出")
			for key, c := range ConnLoggedMap {
				if c == conn {
					delete(ConnLoggedMap, key)
				}
			}
			break
		}
		m := make(map[string]interface{})
		err = json.Unmarshal(b[:n], &m)
		data := m["data"]
		switch m["type"] {
		case "login":
			HandleLogin(data.(map[string]interface{}), conn)
		case "logged_list":
			HandleLoggedList(conn)
		case "message_send":
			HandleMessageSend(data.(map[string]interface{}))
		case "storage_message":
			HandleStorageMessage(data.(map[string]interface{}), conn)
		case "logout":
			HandleLogout(data.(map[string]interface{}))

		default:
			fmt.Printf("没有该类型服务\n")
		}
	}
}
func HandleLogin(data map[string]interface{}, conn net.Conn) {
	r := make(map[string]interface{})
	r["type"] = "login"
	r["code"] = 403
	r["msg"] = "账号或密码错误"
	for _, user := range model.UserAll {
		if user.Username == data["username"] && user.Password == data["password"] {
			if ConnLoggedMap[user.Username] != nil {
				r["code"] = 403
				r["msg"] = "账号已在其它地方登录"
				break
			}
			ConnLoggedMap[user.Username] = conn
			bytesTmp, _ := json.Marshal(user)
			json.Unmarshal(bytesTmp, &data)
			r["code"] = 200
			r["data"] = data
			b, err := json.Marshal(r)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			conn.Write(b)
			break
		}
	}
	b, _ := json.Marshal(r)
	conn.Write(b)
}
func HandleLogout(data map[string]interface{}) {
	delete(ConnLoggedMap, data["username"].(string))
}
func HandleLoggedList(conn net.Conn) {
	r := make(map[string]interface{})
	s := make([]string, 0, 128)
	for key, _ := range ConnLoggedMap {
		s = append(s, key)
	}
	r["data"] = s
	r["type"] = "logged_list"
	b, _ := json.Marshal(r)
	conn.Write(b)
}

func HandleMessageSend(data map[string]interface{}) {
	r := make(map[string]interface{})
	r["type"] = "message_receive"
	r["data"] = data
	fmt.Println(data)
	receiver := data["receiver"].(string)
	var flag bool

	switch data["type"] {
	case "private":
		flag = true
		for username, conn := range ConnLoggedMap {
			if username == receiver {
				flag = false
				r["code"] = 200
				b, _ := json.Marshal(r)
				conn.Write(b)
				fmt.Println("私聊消息发送成功")
				return
			}
		}
		if flag {
			fmt.Println("用户未上线，私聊消息等待发送")
			StorageMessageMap[receiver] = append(StorageMessageMap[receiver], data)
		}
	case "group":
		for _, user := range model.GroupMap[receiver] {
			flag = true
			for name, conn := range ConnLoggedMap {
				if user.Username == name {
					flag = false
					r["code"] = 200
					b, _ := json.Marshal(r)
					conn.Write(b)
					break
				}
			}
			if flag {
				fmt.Println("用户未上线，群聊消息等待发送")
				StorageMessageMap[receiver] = append(StorageMessageMap[receiver], data)
			}

		}
	case "broadcast":
		for _, user := range model.UserAll {
			flag = true
			for name, conn := range ConnLoggedMap {
				if user.Username == name {
					flag = false
					r["code"] = 200
					b, _ := json.Marshal(r)
					conn.Write(b)
					break
				}
			}
			if flag {
				fmt.Println("用户未上线，广播消息等待发送")
				StorageMessageMap[receiver] = append(StorageMessageMap[receiver], data)
			}

		}
	default:
		fmt.Println("没有该消息类型")
	}
}

func HandleStorageMessage(data map[string]interface{}, conn net.Conn) {
	receiver := data["receiver"].(string)
	r := make(map[string]interface{})
	r["type"] = "message_receive"
	for i, message := range StorageMessageMap[receiver] {
		r["code"] = 200
		r["data"] = message
		b, _ := json.Marshal(r)
		conn.Write(b)
		StorageMessageMap[receiver] = append(StorageMessageMap[receiver][:i], StorageMessageMap[receiver][i+1:]...)
	}

}
