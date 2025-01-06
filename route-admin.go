package main

import "net/http"

func registerAdmin() {
	http.HandleFunc("/admin/", adminHandler)
}

func adminHandler(writer http.ResponseWriter, request *http.Request) {
	sendError(writer, http.StatusForbidden, "You are not an admin")
	//	TODO build admin page
	//	 should include data about:
	//	 	-  how many instances have been used,
	//	 	-  how many are left,
	//	 	-  how use using which instance,
	//	 	-  the log file this app has generated thus far
}
