package server

import (
	"taleteller/app"
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
	serviceConfig := app.InitServiceConfig()

	storyStore := store.NewStoryStore(appDB)
	generatorUtils := utils.NewGeneratorUtils()

	storyService := story.NewService(storyStore, serviceConfig.GetPythonServerBaseURL(), generatorUtils)

	dependencies = Dependencies{StoryService: storyService}
	return
}
