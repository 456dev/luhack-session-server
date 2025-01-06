package main

import "net/http"

func registerApp(jwtSecret string) {

	appHandler := func(writer http.ResponseWriter, request *http.Request) {
		userJwt, ok := verifyJwtCookie(writer, request, jwtSecret)
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

	http.HandleFunc("/app/", appHandler)
}
