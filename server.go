package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

import "github.com/golang-jwt/jwt/v5"

type UserJwt struct {
	Admin    bool
	Provider string
	Username string
}

var userInstance map[string]string

var htmlTemplates map[string]*template.Template

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

	err = http.ListenAndServe(serverHost, nil)
	if err != nil {
		panic(err)
		return
	}
}

func authHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/auth/login" {
		hostPath := fmt.Sprintf("https://auth.luhack.uk/?redirect=%s://%s/auth/authenticated", serverProtocol, serverDomain)
		http.Redirect(writer, request, hostPath, http.StatusTemporaryRedirect)
	} else if request.URL.Path == "/auth/authenticated" {
		//	 get jwt param from request
		jwtToken := request.URL.Query().Get("jwt")
		if jwtToken == "" {
			sendError(writer, http.StatusBadRequest, "No jwt token")
			return
		}
		//	 verify jwt, if not valid, return error
		valid, _, err := verifyJwt(jwtToken)
		if err != nil {
			sendError(writer, http.StatusInternalServerError, err.Error())
			return
		}
		if !valid {
			sendError(writer, http.StatusUnauthorized, "Please log in")
			return
		}
		//	 set cookie and redirect to /app/
		http.SetCookie(writer, &http.Cookie{
			Name:  "SessionLogin",
			Value: jwtToken,
			Path:  "/",
		})
		http.Redirect(writer, request, "/app/", http.StatusTemporaryRedirect)
	} else if request.URL.Path == "/auth/logout" {
		http.SetCookie(writer, &http.Cookie{
			Name:     "SessionLogin",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			SameSite: http.SameSiteStrictMode,
		})
		http.Redirect(writer, request, "/auth/login", http.StatusTemporaryRedirect)
	}
}

func appHandler(writer http.ResponseWriter, request *http.Request) {
	userJwt, ok := verifyJwtCookie(writer, request)
	if !ok {
		return
	}
	_, err := writer.Write([]byte("Hi, " + userJwt.Username))
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}
}

func quizHandler(writer http.ResponseWriter, request *http.Request) {
	sendError(writer, http.StatusForbidden, "You are not allowed to access this page")
}

func adminHandler(writer http.ResponseWriter, request *http.Request) {
	sendError(writer, http.StatusForbidden, "You are not an admin")
}

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
				req.URL.Path = "/static" + req.URL.Path
				return
			}
			//	TODO
			//	verify user cookie
			// check if user has an instance set in ProxyPath cookie, if yes, and it matches, proxy to serverHost/proxy/path
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

func verifyJwt(tokenString string) (bool, UserJwt, error) {
	if tokenString == "" {
		log.Println("No token")
		return false, UserJwt{}, nil
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(jwtSecret), nil
	})
	if err != nil {
		log.Fatal(err)
		return false, UserJwt{}, err
	}

	if !token.Valid {
		log.Println("Invalid token")
		return false, UserJwt{}, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claims["exp"].(float64) < float64(time.Now().Unix()) {
			log.Println("Token expired")
			return false, UserJwt{}, nil
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return true, UserJwt{
			Admin:    claims["admin"].(bool),
			Provider: claims["provider"].(string),
			Username: claims["username"].(string),
		}, nil
	} else {
		log.Println("Invalid claims")
		return false, UserJwt{}, err
	}
}

func verifyJwtCookie(writer http.ResponseWriter, request *http.Request) (UserJwt, bool) {
	jwtCookie, err := request.Cookie("SessionLogin")
	if errors.Is(err, http.ErrNoCookie) {
		http.Redirect(writer, request, "/auth/login", http.StatusTemporaryRedirect)
		return UserJwt{}, false
	}
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return UserJwt{}, false
	}
	valid, userJwt, err := verifyJwt(jwtCookie.Value)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return UserJwt{}, false
	}
	if !valid {
		http.Redirect(writer, request, "/auth/login", http.StatusTemporaryRedirect)
		return UserJwt{}, false
	}
	return userJwt, true
}

func sendError(writer http.ResponseWriter, status int, long string) {
	var short string
	switch status {
	case http.StatusNotFound:
		short = "Not Found"
	case http.StatusUnauthorized:
		short = "Please log in"
		long = "<a href=\"/auth/login\">Log in</a>"
	case http.StatusForbidden:
		short = "Forbidden"
	default:
		short = "An unexpected error occurred"
	}

	writer.WriteHeader(status)

	err := htmlTemplates["error.html"].Execute(writer, struct {
		Short string
		Long  string
	}{
		short,
		long,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func parseTemplates() {
	htmlTemplates = make(map[string]*template.Template)
	templates := make([]string, 0)
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".html" {
			templates = append(templates, info.Name())
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, templateName := range templates {
		tmpl, err := template.ParseFiles("templates/" + templateName)
		if err != nil {
			panic(err)
		}
		htmlTemplates[templateName] = tmpl
	}
}
