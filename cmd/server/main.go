package main

import (
	"communication_/internal/server"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		fmt.Printf("err:%v", err)

	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("err:%v", err)
			continue
		}
		go server.HandleConn(conn)
	}
}
