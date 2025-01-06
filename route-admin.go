package main

import "net/http"

func registerAdmin() {
	http.HandleFunc("/admin/", adminHandler)
}

func adminHandler(writer http.ResponseWriter, request *http.Request) {
	sendError(writer, http.StatusForbidden, "You are not an admin")
	//	TODO build admin page
}
