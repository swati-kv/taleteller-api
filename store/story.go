package store

import (
	"context"
	"time"
)

type CreateStoryRequest struct {
	StoryID     string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Mood        string    `db:"mood"`
	Category    string    `db:"category"`
	CustomerID  string    `db:"customer_id"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateSceneRequest struct {
	SceneID     string
	Status      string
	StoryID     string
	SceneNumber int64
}

type InsertImage struct {
	ID        string
	ImagePath string
	SceneID   string
}

type StoryStorer interface {
	Create(ctx context.Context, createRequest CreateStoryRequest) (err error)
	CreateScene(ctx context.Context, request CreateSceneRequest) (err error)
	InsertImage(ctx context.Context, request InsertImage) (err error)
}
