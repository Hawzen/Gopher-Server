package main

import (
	"fmt"
	"net"
)

const (
	// General constants
	KILOBYTE = 1024
	// Connection constants
	CONN_PORT             = "9121"
	CONN_TYPE             = "tcp"
	RECIEVING_BUFFER_SIZE = KILOBYTE * 4
)

func main() {
	listener, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port " + CONN_PORT + " ...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go HandleConn(conn)
	}
}
