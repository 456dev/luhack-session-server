package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			panic(err)
		}
	}(logFile)
	log.SetOutput(logFile)

	var config *Config
	err = parseConfig("config.yml", &config)
	if err != nil {
		panic(err)
	}

	var backendMap *BackendMap
	err = parseBackendMap(config.Session.BackendMap, &backendMap)
	if err != nil {
		panic(err)
	}

	parseTemplates()

	userInstances, allInstances := loadInstances(*backendMap)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	registerError()
	registerAuth(config.Security.JwtSecret, config.Server.Protocol, config.Server.Domain)
	registerAdmin(config.Security.JwtSecret, &userInstances, &allInstances)
	registerApp(config.Security.JwtSecret, config.Session.Title, backendMap.Layout)
	registerProxy(backendMap.LbEndpoint, config.Server.Host, config.Security.JwtSecret, *backendMap, &userInstances, &allInstances)

	registerRoot(config.Server.Host)

	log.Println("Listening on", config.Server.Host)

	err = http.ListenAndServe(config.Server.Host, nil)
	if err != nil {
		panic(err)
		return
	}
}
