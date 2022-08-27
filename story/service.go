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
	GetStoryStatus(ctx context.Context, storyID string) (status string, err error)
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

	req := store.CreateStoryRequest{
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

func (s *service) GetStoryStatus(ctx context.Context, storyID string) (status string, err error) {
	config, err := s.store.GetStoryByID(ctx, storyID)
	if err != nil {
		logger.Errorw(ctx, "error getting story by story ID", "error", err.Error())
		return
	}
	status = config.Status
	return
}
