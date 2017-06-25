package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type authHandler struct {
	router *mux.Router
}

func newAuthHandler(router *mux.Router) *authHandler {
	return &authHandler{router}
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
