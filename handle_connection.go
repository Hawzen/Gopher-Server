package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	COLON         = ":"
	SPACE         = " "
	EMPTY         = ""
	CR            = "\r"
	LF            = "\n"
	CRLF          = CR + LF
	CR_BYTES      = byte('\r')
	LF_BYTES      = byte('\n')
	PROTOCOL_USED = "HTTP/1.1"
)

type URI string

// SUPPORTED_METHODS
var SUPPORTED_METHODS = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}
var SUPPORTED_PROTOCOLS = []string{"HTTP/1.1"}

type Request struct {
	method   string
	uri      URI
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

func HandleConn(conn net.Conn) {
	// TODO: Establish deadline for reading/writing
	defer conn.Close()

	// Reading the request
	request, err := handle_conn_read(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("RECIEVED" + COLON + LF + request.stringify())

	// Writing the response
	response, err := handle_conn_write(conn, request)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("SENDING" + COLON + LF + response.stringify())

	fmt.Println("______________________________")
}

func handle_conn_read(conn net.Conn) (Request, error) {
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
	request_method := preprocess_http_text(request_line_tokens[0])
	request_uri := preprocess_http_text(request_line_tokens[1])
	request_protocol := preprocess_http_text(request_line_tokens[2])

	if !slices.Contains(SUPPORTED_PROTOCOLS, request_protocol) {
		return Request{}, fmt.Errorf("Invalid request protocol: %s", request_protocol)
	}
	// Check if URI is valid
	if !slices.Contains(SUPPORTED_METHODS, request_method) {
		return Request{}, fmt.Errorf("Invalid request method: %s", request_method)
	}

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
		uri:      URI(request_uri),
		protocol: request_protocol,
		headers:  request_headers,
		body:     request_body,
	}

	return request, nil
}

func handle_conn_write(conn net.Conn, request Request) (Response, error) {
	var response Response
	if PathToHandler[request.uri] == nil {
		response = PathToHandler[SPECIAL_PATH_ELSE](request)
	} else {
		response = PathToHandler[request.uri](request)
	}

	fill_boring_response_fields(&response)

	serialized_response := response.stringify()

	_, err := conn.Write([]byte(serialized_response))
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func (request Request) stringify() string {
	var request_builder strings.Builder
	request_builder.WriteString(request.method)
	request_builder.WriteString(SPACE)
	request_builder.WriteString(string(request.uri))
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

func preprocess_http_text(text string) string {
	return strings.TrimSpace(
		strings.ToUpper(
			text,
		),
	)
}

func fill_boring_response_fields(response *Response) {
	if response.headers == nil {
		response.headers = make(map[string]string)
	}
	response.headers["Content-Type"] = "text/plain; charset=utf-8"
	response.headers["Content-Length"] = fmt.Sprintf("%d", len(response.body))
	response.protocol = PROTOCOL_USED
}
