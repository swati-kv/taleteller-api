package store

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type storyStore struct {
	db *sqlx.DB
}

func (s storyStore) Create(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func NewStoryStore(db *sqlx.DB) StoryStorer {
	return &storyStore{
		db: db,
	}
}
