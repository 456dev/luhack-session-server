package main

import (
	"net/http"
	"os"
)

var userInstance map[string]string

var serverHost string
var serverDomain string
var serverProtocol string
var jwtSecret string

func main() {
	// TODO add proper logging
	args := os.Args
	backendMapPath := "backend-map.yaml"
	if len(args) == 2 {
		backendMapPath = args[1]
	}

	var backendMap *BackendMap
	err := parseBackendMap(backendMapPath, &backendMap)
	if err != nil {
		panic(err)
	}

	// TODO don't hardcode these
	serverHost = "localhost:8080"
	// TODO make the default session.luhack.uk
	serverDomain = serverHost
	serverProtocol = "http"
	jwtSecret = "yWGOSeOmQu5RG2m8Wgz4KO2kZmD4Yoz5XdNz5sGS4_E"

	parseTemplates()

	http.HandleFunc("/auth/", authHandler)
	http.HandleFunc("/app/", appHandler)
	http.HandleFunc("/quiz/", quizHandler)
	http.HandleFunc("/admin/", adminHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	registerRoot(serverHost)

	registerProxy(backendMap.LbEndpoint)

	err = http.ListenAndServe(serverHost, nil)
	if err != nil {
		panic(err)
		return
	}
}
