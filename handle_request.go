package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	SPECIAL_PATH_ELSE = "else"
)

var current_directory, _ = os.Getwd()

func index(request Request) Response {
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		return internal_error(request)
	}

	var index_page_builder strings.Builder
	index_page_builder.WriteString("<h1>" + current_directory + " /</h1>")
	index_page_builder.WriteString("<hr>\n")
	for _, file := range files {
		file_name := file.Name()
		file_path := "/" + file_name
		index_page_builder.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a><br/>", file_path, file_name))
	}
	body := index_page_builder.String()

	return Response{
		status_code:    200,
		status_message: "OK",
		headers: map[string]string{
			"Content-Type": "text/html",
		},
		body: body,
	}
}

func not_found(request Request) Response {
	return Response{
		status_code:    404,
		status_message: "Not Found",
		body:           "Not Found, homie",
	}
}

func internal_error(request Request) Response {
	return Response{
		status_code:    500,
		status_message: "Internal Server Error",
		body:           "Internal Server Error, homie",
	}
}

func serve_static_file(request Request) Response {
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
	"/":      index,
	"/index": index,
	// Special cases
	SPECIAL_PATH_ELSE: serve_static_file,
}
