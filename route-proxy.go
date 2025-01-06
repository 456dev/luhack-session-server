package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"strings"
)

func registerProxy(backendTarget string, serverHost string, jwtSecret string) {
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = serverHost
			req.URL.Path = "/error"
			req.URL.RawQuery = "code=403&message=Failed to authenticate"

			jwtCookie, err := req.Cookie("SessionLogin")
			if err != nil {
				return
			}
			valid, _, err := verifyJwt(jwtCookie.Value, jwtSecret)
			if err != nil || !valid {
				return
			}

			req.URL.Scheme = "http"
			req.URL.Host = backendTarget
			req.URL.Path = "/" + strings.TrimPrefix(req.URL.Path, "/proxy")

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
