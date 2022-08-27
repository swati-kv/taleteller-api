package store

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type storyStore struct {
	db *sqlx.DB
}

func (s storyStore) InsertImage(ctx context.Context, request InsertImage) (err error) {
	_, err = s.db.ExecContext(ctx, insertImage,
		request.ID,
		request.ImagePath,
		request.SceneID,
		false,
		time.Now(),
		time.Now(),
	)
	return
}

func (s storyStore) Create(ctx context.Context, c CreateStoryRequest) (err error) {
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

func (s storyStore) CreateScene(ctx context.Context, request CreateSceneRequest) (err error) {
	//TODO implement me
	_, err = s.db.ExecContext(ctx, createScene,
		request.SceneID,
		request.StoryID,
		request.Status,
		request.SceneNumber,
		time.Now(),
		time.Now(),
	)
	return
}

func NewStoryStore(db *sqlx.DB) StoryStorer {
	return &storyStore{
		db: db,
	}
}
