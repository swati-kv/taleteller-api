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

	router.Handle("/story/{id}/scene", story.HandleCreateScene(dependencies.StoryService)).Methods(http.MethodPost)

	router.Handle("/stories/{id}",
		story.HandleGetStory(dependencies.StoryService),
	).Methods(http.MethodGet)

	router.Handle("/stories", story.HandleListStories(dependencies.StoryService)).Methods(http.MethodGet)

	router.Handle("/stories/{id}/publish", story.HandlePublishStory(dependencies.StoryService)).Methods(http.MethodPost)

	return
}
