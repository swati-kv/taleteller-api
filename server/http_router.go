package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"taleteller/story"
)

func initRouter(dependencies Dependencies) (router *mux.Router) {
	router = mux.NewRouter()
	router.StrictSlash(true)

	router.Handle("/stories", story.HandleStoryCreate(dependencies.StoryService)).Methods(http.MethodPost)

	router.Handle("/stories/{id}",
		story.HandleGetStoryStatus(dependencies.StoryService),
	).Methods(http.MethodGet)
	router.Handle("/stories", story.HandleListStories(dependencies.StoryService)).Methods(http.MethodGet)

	return
}
