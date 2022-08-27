package story

import (
	"context"
	"taleteller/logger"
	"taleteller/store"
	"taleteller/utils"
	"time"
)

type Service interface {
	Create(ctx context.Context, createRequest CreateStoryRequest) (err error)
	GetStory(ctx context.Context, storyID string) (storyDetails store.Story, err error)
	List(ctx context.Context, status string) (stories []store.Story, err error)
	UpdateScene(ctx context.Context, storyID string, sceneID string, selectedImage string) (updatedScene store.Scene, err error)
}

type service struct {
	store          store.StoryStorer
	generatorUtils utils.IDGeneratorUtils
}

func NewService(store store.StoryStorer, generatorUtils utils.IDGeneratorUtils) *service {
	return &service{
		store:          store,
		generatorUtils: generatorUtils,
	}
}

func (s *service) Create(ctx context.Context, createRequest CreateStoryRequest) (err error) {

	storyID, err := s.generatorUtils.GenerateIDWithPrefix("sto_")
	if err != nil {
		logger.Error(ctx, "error generating ID", err.Error())
		return
	}

	req := store.Story{
		StoryID:     storyID,
		Name:        createRequest.Name,
		Description: createRequest.Description,
		Mood:        createRequest.Mood,
		Category:    createRequest.Category,
		CustomerID:  "cus_123",
		Status:      "processing",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	err = s.store.Create(ctx, req)
	if err != nil {
		logger.Error(ctx, "error creating story", err.Error())
		return
	}
	return
}

func (s *service) GetStory(ctx context.Context, storyID string) (storyDetails store.Story, err error) {
	storyDetails, err = s.store.GetStoryByID(ctx, storyID)
	if err != nil {
		logger.Errorw(ctx, "error getting story by story ID", "error", err.Error())
		return
	}
	return
}

func (s *service) List(ctx context.Context, status string) (stories []store.Story, err error) {
	stories, err = s.store.List(ctx, status)
	if err != nil {
		logger.Error(ctx, "error getting stories", err.Error())
		return
	}
	return
}

func (s *service) UpdateScene(ctx context.Context, storyID string, sceneID string, selectedImage string) (updatedScene store.Scene, err error) {

	updatedScene, err = s.store.UpdateScene(ctx, storyID, sceneID, selectedImage)
	if err != nil {
		logger.Error(ctx, "error updating scene", err.Error())
		return
	}
	return
}
