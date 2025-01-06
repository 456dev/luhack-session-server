package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"strings"
)

func registerProxy(backendTarget string) {
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = backendTarget
			req.URL.Path = "/" + strings.TrimPrefix(req.URL.Path, "/proxy")

			//	TODO verify user cookie
		},
		ModifyResponse: func(response *http.Response) error {
			return nil
		},
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true, // Disable keep-alives to ensure new connections for each request
		},
	}

	http.Handle("/proxy/", proxy)
}
