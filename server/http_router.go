package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"taleteller/service"
)

func initRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.StrictSlash(true)

	router.Handle("/sample", service.HandlerSample()).Methods(http.MethodGet)

	return
}
