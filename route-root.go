package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
)

func registerRoot(serverHost string) {
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			if req.URL.Path == "/" {
				req.URL.Path = "/app/"
				req.URL.Host = serverHost
				req.URL.Scheme = "http"

				return
			} else if req.URL.Path == "/favicon.ico" || req.URL.Path == "/robots.txt" {
				req.URL.Host = serverHost
				req.URL.Scheme = "http"
				req.URL.Path = "/_session_server/static" + req.URL.Path
				return
			}

			req.URL.Host = serverHost
			req.URL.Scheme = "http"
			req.URL.Path = "/proxy" + req.URL.Path
			return
		},
		ModifyResponse: func(response *http.Response) error {
			return nil
		},
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
	}
	http.Handle("/", proxy)
}
