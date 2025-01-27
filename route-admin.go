package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

func getLastNLines(logContent string, n int) string {
	lines := strings.Split(logContent, "\n")
	if len(lines) <= n {
		return logContent
	}
	return strings.Join(lines[len(lines)-n:], "\n")
}

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

		if strings.HasPrefix(request.URL.Path, "/admin/release/") {
			uid := UID(request.URL.Path[len("/admin/release/"):])
			err := uid.releaseInstance(userInstances, allInstances)
			if err != nil {
				sendError(writer, http.StatusInternalServerError, err.Error())
				return
			}
			http.Redirect(writer, request, "/admin/", http.StatusTemporaryRedirect)
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

		pageData.Logs = getLastNLines(string(logFile), 1000)

		pageData.TotalInstances = len(*allInstances)
		pageData.UsedInstances = 0
		pageData.FreeInstances = 0
		pageData.UserInstances = *userInstances

		for _, instance := range *allInstances {
			if !instance {
				pageData.UsedInstances++
			}
			if instance {
				pageData.FreeInstances++
			}
		}

		err = htmlTemplates["admin.html"].Execute(writer, pageData)
		if err != nil {
			log.Println("Failed to execute admin template:", err)
			sendError(writer, http.StatusInternalServerError, err.Error())
			return
		}
	}

	http.HandleFunc("/admin/", adminHandler)

}
