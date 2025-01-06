package main

import "net/http"

func adminHandler(writer http.ResponseWriter, request *http.Request) {
	sendError(writer, http.StatusForbidden, "You are not an admin")
	//	TODO build admin page
}
