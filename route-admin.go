package main

import (
	"errors"
	"log"
	"net/http"
	"os"
)

func registerAdmin(jwtSecret string, userInstances *map[UID]Instance, allInstances *map[Instance]bool) {
	adminHandler := func(writer http.ResponseWriter, request *http.Request) {
		sessionCookie, err := request.Cookie("SessionLogin")
		if errors.Is(err, http.ErrNoCookie) {
			sendError(writer, http.StatusUnauthorized, "Please log in")
			return
		}
		if err != nil {
			log.Println("Failed to get session cookie:", err)
			sendError(writer, http.StatusInternalServerError, err.Error())
			return
		}

		valid, user, err := verifyJwt(sessionCookie.Value, jwtSecret)
		if err != nil {
			log.Println("Failed to verify jwt:", err)
			sendError(writer, http.StatusInternalServerError, err.Error())
			return
		}
		if !valid {
			sendError(writer, http.StatusUnauthorized, "Please log in")
			return
		}
		if !user.Admin {
			sendError(writer, http.StatusForbidden, "Forbidden")
			return
		}

		pageData := struct {
			Logs           string
			TotalInstances int
			UsedInstances  int
			FreeInstances  int
			UserInstances  map[UID]Instance
		}{}

		logFile, err := os.ReadFile("log.txt")
		if err != nil {
			log.Println("Failed to read log file:", err)
			sendError(writer, http.StatusInternalServerError, err.Error())
			return
		}

		pageData.Logs = string(logFile)

		pageData.TotalInstances = len(*allInstances)
		pageData.UsedInstances = len(*userInstances)
		pageData.FreeInstances = pageData.TotalInstances - pageData.UsedInstances
		pageData.UserInstances = *userInstances

		err = htmlTemplates["admin.html"].Execute(writer, pageData)
		if err != nil {
			log.Println("Failed to execute admin template:", err)
			sendError(writer, http.StatusInternalServerError, err.Error())
			return
		}
	}

	http.HandleFunc("/admin/", adminHandler)
}
