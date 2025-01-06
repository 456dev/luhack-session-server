package main

import "net/http"

func registerQuiz() {
	http.HandleFunc("/quiz/", quizHandler)
}

func quizHandler(writer http.ResponseWriter, request *http.Request) {
	// TODO get the quiz page content and stuff
	sendError(writer, http.StatusForbidden, "You are not allowed to access this page")
}
