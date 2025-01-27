package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
)

type UserJwt struct {
	Admin    bool
	Provider string
	Username string
}

func registerAuth(jwtSecret string, serverProtocol string, serverDomain string) {

	authHandler := func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/auth/login" || request.URL.Path == "/auth/login/" {
			hostPath := fmt.Sprintf("https://auth.luhack.uk/?redirect=%s://%s/auth/authenticated", serverProtocol, serverDomain)
			http.Redirect(writer, request, hostPath, http.StatusTemporaryRedirect)
		} else if request.URL.Path == "/auth/authenticated" || request.URL.Path == "/auth/authenticated/" {
			//	 get jwt param from request
			jwtToken := request.URL.Query().Get("jwt")
			if jwtToken == "" {
				sendError(writer, http.StatusBadRequest, "No jwt token")
				log.Println("No jwt token from user with IP", request.RemoteAddr, "with user agent", request.UserAgent())
				return
			}
			//	 verify jwt, if not valid, return error
			valid, _, err := verifyJwt(jwtToken, jwtSecret)
			if err != nil {
				sendError(writer, http.StatusInternalServerError, err.Error())
				log.Println("Failed to verify jwt from user with IP", request.RemoteAddr, "with user agent", request.UserAgent())
				return
			}
			if !valid {
				sendError(writer, http.StatusUnauthorized, "Please log in")
				log.Println("Invalid jwt from user with IP", request.RemoteAddr, "with user agent", request.UserAgent())
				return
			}
			//	 set cookie and redirect to /app/
			http.SetCookie(writer, &http.Cookie{
				Name:  "SessionLogin",
				Value: jwtToken,
				Path:  "/",
			})
			http.Redirect(writer, request, "/app/", http.StatusTemporaryRedirect)
		} else if request.URL.Path == "/auth/logout" || request.URL.Path == "/auth/logout/" {
			http.SetCookie(writer, &http.Cookie{
				Name:     "SessionLogin",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				SameSite: http.SameSiteStrictMode,
			})
			err := htmlTemplates["logout.html"].Execute(writer, nil)
			if err != nil {
				sendError(writer, http.StatusInternalServerError, err.Error())
				log.Println("Failed to render logout.html from user with IP", request.RemoteAddr, "with user agent", request.UserAgent())
				return
			}
		}
	}

	http.HandleFunc("/auth/", authHandler)

}

func verifyJwt(tokenString string, jwtSecret string) (bool, UserJwt, error) {
	if tokenString == "" {
		log.Println("No token")
		return false, UserJwt{}, nil
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		if err.Error() == "token has invalid claims: token is expired" {
			return false, UserJwt{}, nil
		}
		log.Println("Failed to parse token:", err)
		return false, UserJwt{}, err
	}

	if !token.Valid {
		log.Println("Invalid token")
		return false, UserJwt{}, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return true, UserJwt{
			Admin:    claims["admin"].(bool),
			Provider: claims["provider"].(string),
			Username: claims["username"].(string),
		}, nil
	} else {
		log.Println("Failed to parse claims")
		return false, UserJwt{}, err
	}
}

func verifyJwtCookie(writer http.ResponseWriter, request *http.Request, jwtSecret string) (UserJwt, bool) {
	jwtCookie, err := request.Cookie("SessionLogin")
	if errors.Is(err, http.ErrNoCookie) {
		http.Redirect(writer, request, "/auth/login", http.StatusTemporaryRedirect)
		return UserJwt{}, false
	}
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		log.Println("Failed to get jwt cookie from user with IP", request.RemoteAddr, "with user agent", request.UserAgent())
		return UserJwt{}, false
	}
	valid, userJwt, err := verifyJwt(jwtCookie.Value, jwtSecret)
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		log.Println("Failed to verify jwt from user with IP", request.RemoteAddr, "with user agent", request.UserAgent())
		return UserJwt{}, false
	}
	if !valid {
		http.Redirect(writer, request, "/auth/login", http.StatusTemporaryRedirect)
		return UserJwt{}, false
	}
	return userJwt, true
}
