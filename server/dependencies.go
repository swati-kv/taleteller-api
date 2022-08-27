package server

import (
	"taleteller/db"
	"taleteller/store"
	"taleteller/story"
	"taleteller/utils"
)

type Dependencies struct {
	StoryService story.Service
}

func NewDependencies() (dependencies Dependencies, err error) {
	appDB := db.Get()

	storyStore := store.NewStoryStore(appDB)
	generatorUtils := utils.NewGeneratorUtils()

	storyService := story.NewService(storyStore, generatorUtils)

	dependencies = Dependencies{StoryService: storyService}
	return
}
