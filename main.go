package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

const (
	// General constants
	KILOBYTE = 1024
	// Connection constants
	CONN_PORT             = "9121"
	CONN_TYPE             = "tcp"
	RECIEVING_BUFFER_SIZE = KILOBYTE * 4
	// HTTP constants
	CLRF     = "\r\n"
	NEW_LINE = "\n"
	COLON    = ":"
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
	// TODO: Establish deadline for reading/writing
	defer conn.Close()

	// Reading the request
	request, err := handleRead(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("RECIEVED" + COLON + NEW_LINE + request)

	// Writing the response
	response, err := handleWrite(conn, request)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("SENDING" + COLON + NEW_LINE + response + NEW_LINE + "_______________________")

}

func handleRead(conn net.Conn) (string, error) {
	// TODO: Validate request
	// TODO: Return the request as a struct
	buf := make([]byte, RECIEVING_BUFFER_SIZE)
	var request bytes.Buffer
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return "", err
		}
		request.Write(buf[:n])
		if strings.HasSuffix(request.String(), CLRF+CLRF) {
			break
		}
	}
	return request.String(), nil
}

func handleWrite(conn net.Conn, message string) (string, error) {
	// TODO: Return the response as a struct
	var response_builder strings.Builder
	status_line := "HTTP/1.1 200 OK" + NEW_LINE
	payload := "Hello World, go! " + NEW_LINE + "رواية البؤساء"

	headers := map[string]string{
		"Content-Type":   "text/plain; charset=utf-8",
		"Content-Length": fmt.Sprintf("%d", len(payload)),
	}

	response_builder.WriteString(status_line)
	for k, v := range headers {
		response_builder.WriteString(k)
		response_builder.WriteString(COLON)
		response_builder.WriteString(v)
		response_builder.WriteString(NEW_LINE)
	}
	response_builder.WriteString(CLRF)
	response_builder.WriteString(payload)

	response := response_builder.String()

	_, err := conn.Write([]byte(response))
	if err != nil {
		return "", err
	}

	return response, nil
}
