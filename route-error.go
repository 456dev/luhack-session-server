package main

import (
	"log"
	"net/http"
	"strconv"
)

func registerError() {
	handleError := func(writer http.ResponseWriter, request *http.Request) {
		code := request.URL.Query().Get("code")
		message := request.URL.Query().Get("message")

		intCode, err := strconv.Atoi(code)
		if err != nil || intCode < 400 || intCode > 599 {
			intCode = http.StatusInternalServerError
			message = "Invalid error code"
		}

		sendError(writer, intCode, message)
	}
	// TODO make this actually work, currently error pages that are proxied return 400 codes

	http.HandleFunc("/error", handleError)
	http.HandleFunc("/error/", handleError)
}

func sendError(writer http.ResponseWriter, status int, long string) {
	var short string
	switch status {
	case http.StatusNotFound:
		short = "Not Found"
	case http.StatusUnauthorized:
		short = "Please log in"
		long = "<a href=\"/auth/login\">Log in</a>"
	case http.StatusForbidden:
		short = "Forbidden"
	default:
		short = "An unexpected error occurred"
	}

	writer.WriteHeader(status)

	err := htmlTemplates["error.html"].Execute(writer, struct {
		Short string
		Long  string
	}{
		short,
		long,
	})
	if err != nil {
		log.Println("Failed to execute error template:", err)
		return
	}
}
