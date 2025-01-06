package main

import (
	"log"
	"net/http"
)

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
		log.Println(err)
		return
	}
}
