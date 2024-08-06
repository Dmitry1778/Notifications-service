package api

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

func headerStr(w http.ResponseWriter, r *http.Request) (string, error) {
	var accessToken string
	header := r.Header.Get(authorizationHeader)
	if header == "" {
		log.Println("empty header")
		return "", errors.New("error")
	}
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		log.Println("Invalid Authorization header format")
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return "", errors.New("error")
	} else {
		accessToken = parts[1]
	}
	return accessToken, nil
}
