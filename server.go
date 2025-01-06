package main

import (
	"net/http"
)

func main() {
	// TODO add proper logging

	var config *Config
	err := parseConfig("config.yml", &config)
	if err != nil {
		panic(err)
	}

	var backendMap *BackendMap
	err = parseBackendMap(config.Session.BackendMap, &backendMap)
	if err != nil {
		panic(err)
	}

	parseTemplates()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	registerError()
	registerAuth(config.Security.JwtSecret, config.Server.Protocol, config.Server.Domain)
	registerAdmin()
	registerApp(config.Security.JwtSecret, config.Session.Title, backendMap.Layout)
	registerProxy(backendMap.LbEndpoint, config.Server.Host, config.Security.JwtSecret, *backendMap)

	registerRoot(config.Server.Host)

	err = http.ListenAndServe(config.Server.Host, nil)
	if err != nil {
		panic(err)
		return
	}
}
