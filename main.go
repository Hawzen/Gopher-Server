package main

import (
	"fmt"
	"net"
)

const (
	CONN_PORT = "9121"
	CONN_TYPE = "tcp"
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
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// Just print whatever you receive
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("Received: %s\n", buf[:n])
	}
}
