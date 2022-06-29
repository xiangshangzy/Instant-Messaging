package main

import (
	"communication_/internal/client"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Printf("%v", err)
	}
	client.HandleConn(conn)
}
