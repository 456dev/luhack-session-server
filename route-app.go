package main

import (
	"log"
	"net/http"
)

func registerApp(jwtSecret string, appTitle string, layout []Layout) {
	type PageData struct {
		Title    string
		Username string
		Layout   []Layout
	}

	appHandler := func(writer http.ResponseWriter, request *http.Request) {
		userJwt, ok := verifyJwtCookie(writer, request, jwtSecret)
		if !ok {
			sendError(writer, http.StatusUnauthorized, "Please log in")
			return
		}

		data := PageData{
			Title:    appTitle,
			Username: userJwt.Username,
			Layout:   layout,
		}

		err := htmlTemplates["app.html"].Execute(writer, data)
		if err != nil {
			log.Println(err)
			sendError(writer, http.StatusInternalServerError, "An unexpected error occurred")
			return
		}

		log.Println("User", buildUid(userJwt), "accessed app")
	}

	http.HandleFunc("/app/", appHandler)
}
