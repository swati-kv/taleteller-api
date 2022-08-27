package story

import (
	"context"
	"taleteller/store"
)

type Service interface {
	Create(ctx context.Context, createRequest CreateStoryRequest) (err error)
}

type service struct {
	store store.StoryStorer
}

func NewService(store store.StoryStorer) Service {
	return &service{
		store: store,
	}
}

func (s *service) Create(ctx context.Context, createRequest CreateStoryRequest) (err error) {
	return
}
