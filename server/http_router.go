package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"taleteller/scene"
	"taleteller/story"
)

func initRouter(dependencies Dependencies) (router *mux.Router) {
	router = mux.NewRouter()
	router.StrictSlash(true)

	router.Handle("/sample", scene.HandlerSample()).Methods(http.MethodGet)

	router.Handle("/story/{id}/scene", story.HandleCreateScene()).Methods(http.MethodPost)

	return
}
