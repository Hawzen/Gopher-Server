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
	COLON    = ":"
	SPACE    = " "
	EMPTY    = ""
	CR       = "\r"
	LF       = "\n"
	CRLF     = CR + LF
	CR_BYTES = byte('\r')
	LF_BYTES = byte('\n')
)

type Request struct {
	method   string
	uri      string
	protocol string
	headers  map[string]string
	body     string
}

type Response struct {
	protocol       string
	status_code    int
	status_message string
	headers        map[string]string
	body           string
}

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
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	// TODO: Establish deadline for reading/writing
	defer conn.Close()

	// Reading the request
	request, err := handleRead(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("RECIEVED" + COLON + LF + request.stringify())

	// Writing the response
	response, err := handleWrite(conn, request)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("SENDING" + COLON + LF + response.stringify())

	fmt.Println("______________________________")
}

func handleRead(conn net.Conn) (Request, error) {
	buf := make([]byte, RECIEVING_BUFFER_SIZE)
	var bytes_request bytes.Buffer
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return Request{}, err
		}
		bytes_request.Write(buf[:n])
		if n >= 2 && buf[n-2] == CR_BYTES && buf[n-1] == LF_BYTES {
			break
		}
	}

	serialized_request := bytes_request.String()

	// Parse and validate the request
	request_split_by_line := strings.Split(serialized_request, LF)
	if len(request_split_by_line) < 2 {
		return Request{}, fmt.Errorf("Invalid request: request lines less than 2")
	}

	request_line := request_split_by_line[0]
	request_line_tokens := strings.Split(request_line, SPACE)
	if len(request_line_tokens) != 3 {
		return Request{}, fmt.Errorf("Invalid request line: %s", request_line)
	}
	request_method := request_line_tokens[0]
	request_uri := request_line_tokens[1]
	request_protocol := request_line_tokens[2]

	// Parse the headers (if any)
	CRLF_line_index := -1
	request_headers := make(map[string]string)
	for i, header_line := range request_split_by_line[1:] {
		if header_line == CR || header_line == EMPTY {
			CRLF_line_index = i
			break
		}
		header_line_tokens := strings.Split(header_line, COLON)
		if len(header_line_tokens) < 2 {
			// Print hex of header_line
			fmt.Printf("%x\n", header_line)
			return Request{}, fmt.Errorf("Invalid header line: %s, %d", header_line, len(header_line_tokens))
		}
		request_headers[header_line_tokens[0]] = strings.Join(header_line_tokens[1:], COLON)
	}

	if CRLF_line_index == -1 {
		return Request{}, fmt.Errorf("Invalid request: no CRLF found")
	}

	// Parse the body
	request_body := strings.Join(request_split_by_line[CRLF_line_index+1:], LF)

	request := Request{
		method:   request_method,
		uri:      request_uri,
		protocol: request_protocol,
		headers:  request_headers,
		body:     request_body,
	}

	return request, nil
}

func handleWrite(conn net.Conn, request Request) (Response, error) {
	// TODO: Return the response as a struct
	response_body := handleRequest(request)

	response := Response{
		protocol:       "HTTP/1.1",
		status_code:    200,
		status_message: "OK",
		headers: map[string]string{
			"Content-Type":   "text/plain; charset=utf-8",
			"Content-Length": fmt.Sprintf("%d", len(response_body)),
		},
		body: response_body,
	}

	// Building the response
	serialized_response := response.stringify()

	_, err := conn.Write([]byte(serialized_response))
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func handleRequest(request Request) string {
	// Stub for now
	return `
	Hello World, go! 
	رواية البؤساء
	`
}

// Methods for stringifying the request and response
func (request Request) stringify() string {
	var request_builder strings.Builder
	request_builder.WriteString(request.method)
	request_builder.WriteString(SPACE)
	request_builder.WriteString(request.uri)
	request_builder.WriteString(SPACE)
	request_builder.WriteString(request.protocol)
	request_builder.WriteString(LF)

	for k, v := range request.headers {
		request_builder.WriteString(k)
		request_builder.WriteString(COLON)
		request_builder.WriteString(v)
		request_builder.WriteString(LF)
	}
	request_builder.WriteString(CRLF)
	request_builder.WriteString(request.body)

	return request_builder.String()
}

func (response Response) stringify() string {
	var response_builder strings.Builder
	response_builder.WriteString(response.protocol)
	response_builder.WriteString(SPACE)
	response_builder.WriteString(fmt.Sprintf("%d", response.status_code))
	response_builder.WriteString(SPACE)
	response_builder.WriteString(response.status_message)
	response_builder.WriteString(LF)

	for k, v := range response.headers {
		response_builder.WriteString(k)
		response_builder.WriteString(COLON)
		response_builder.WriteString(v)
		response_builder.WriteString(LF)
	}
	response_builder.WriteString(CRLF)
	response_builder.WriteString(response.body)

	return response_builder.String()
}
