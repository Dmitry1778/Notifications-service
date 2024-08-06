package api

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"notify/token_generator"
	"strconv"
	"strings"
)

func (a *Api) getValues(w http.ResponseWriter, r *http.Request) (*User, error) {
	empIDstr := chi.URLParam(r, "publisherID")
	empIDstr = strings.Trim(empIDstr, "{}")
	publisherID, err := strconv.Atoi(empIDstr)
	if err != nil {
		log.Printf("strconv.Atoi: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	accessToken, err := headerStr(w, r)
	if err != nil {
		log.Printf("Error get header string: %v", err)
		return nil, err
	}
	token := token_generator.New([]byte("secret-key"), a.db)
	subscribeID, err := token.ParseToken(accessToken)
	//userID, err := a.auth.ParseToken(accessToken)
	if err != nil {
		log.Printf("Error parse token: %v", err)
		return nil, nil
	}
	return &User{subscribeID, publisherID}, nil
}

type User struct {
	subscribeID int
	publisherID int
}
