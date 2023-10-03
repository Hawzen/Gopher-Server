package main

import (
	"fmt"
	"os"
)

const (
	SPECIAL_PATH_ELSE = "else"
)

func dummy_response(request Request) Response {
	body := `
	Hello World, go! 
	رواية البؤساء
	`
	return Response{
		status_code:    200,
		status_message: "OK",
		body:           body,
	}
}

func not_found(request Request) Response {
	return Response{
		status_code:    404,
		status_message: "Not Found",
		body:           "Not Found, homie",
	}
}

func serve_static_file(request Request) Response {
	current_directory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return not_found(request)
	}

	file_content, err := os.ReadFile(current_directory + string(request.uri))
	if err != nil {
		fmt.Println(err)
		return not_found(request)
	}

	return Response{
		status_code:    200,
		status_message: "OK",
		body:           string(file_content),
	}

}

var PathToHandler = map[URI]func(Request) Response{
	"/":           dummy_response,
	"/index.html": dummy_response,
	// Special cases
	SPECIAL_PATH_ELSE: serve_static_file,
}
