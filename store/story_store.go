package store

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type storyStore struct {
	db *sqlx.DB
}

func (s storyStore) List(ctx context.Context) (stories []Story, err error) {
	err = s.db.SelectContext(ctx, &stories, getStories)
	return
}

func (s storyStore) Create(ctx context.Context, c Story) (err error) {
	//TODO implement me
	_, err = s.db.ExecContext(ctx, createStory,
		c.StoryID,
		c.Name,
		c.Mood,
		c.Category,
		c.Description,
		c.CustomerID,
		c.Status,
		time.Now(),
		time.Now())
	return
}
func (s storyStore) GetStoryByID(ctx context.Context, storyID string) (storyResponse GetStoryResponse, err error) {
	err = s.db.GetContext(ctx, &storyResponse, getStoryByID, storyID)
	return
}
func NewStoryStore(db *sqlx.DB) *storyStore {
	return &storyStore{
		db: db,
	}
}
