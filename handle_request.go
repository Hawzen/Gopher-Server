package main

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

var PathToHandler = map[URI]func(Request) Response{
	"/":            dummy_response,
	"/index.html":  dummy_response,
	"/favicon.ico": not_found,
	// Special cases
	SPECIAL_PATH_ELSE: not_found,
}
