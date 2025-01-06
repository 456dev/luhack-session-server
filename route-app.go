package main

import "net/http"

func appHandler(writer http.ResponseWriter, request *http.Request) {
	userJwt, ok := verifyJwtCookie(writer, request)
	if !ok {
		return
	}
	_, err := writer.Write([]byte("Hi, " + userJwt.Username))
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	//	TODO build app page
}
