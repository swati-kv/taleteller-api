package server

import (
	"taleteller/db"
	"taleteller/store"
	"taleteller/story"
)

type Dependencies struct {
	StoryService story.Service
}

func NewDependencies() (dependencies Dependencies, err error) {
	appDB := db.Get()

	storyStore := store.NewStoryStore(appDB)

	storyService := story.NewService(storyStore)

	dependencies = Dependencies{StoryService: storyService}
	return
}
