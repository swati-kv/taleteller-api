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
type Scene struct {
	ID                  string     `db:"id,omitempty" json:"id,omitempty"`
	StoryID             string     `db:"story_id,omitempty" json:"story-id,omitempty"`
	GeneratedAudioID    string     `db:"generated_audio_id,omitempty" json:"generated-audio-id,omitempty"`
	GeneratedAudioPath  string     `json:"generated-audio-path" db:"path,omitempty"`
	BackgroundAudioPath string     `db:"background_audio_path,omitempty" json:"background-audio-path,omitempty"`
	Status              string     `db:"status,omitempty" json:"status,omitempty"`
	SceneNumber         int        `db:"scene_number,omitempty" json:"scene-number,omitempty"`
	SelectedImage       string     `db:"selected_image" json:"selected-image"`
	SelectedImagePath   string     `db:"path"`
	CreatedAt           *time.Time `db:"created_at" json:"created-at,omitempty"`
	UpdatedAt           *time.Time `db:"updated_at" json:"updated-at,omitempty"`
}

type GeneratedAudio struct {
	ID   string `db:"id"`
	Path string `db:"path"`
}

type Story struct {
	StoryID      string    `db:"id" json:"story_id,omitempty"`
	Name         string    `db:"name" json:"name,omitempty"`
	Description  string    `db:"description" json:"description,omitempty"`
	Mood         string    `db:"mood" json:"mood,omitempty"`
	Category     string    `db:"category" json:"category,omitempty"`
	CustomerID   string    `db:"customer_id" json:"customer_id,omitempty"`
	Status       string    `db:"status" json:"status,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	SceneDetails []Scene   `json:"scene-details,omitempty"`
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
	GetStoryByID(ctx context.Context, storyID string) (storyDetails Story, err error)
	Create(ctx context.Context, createRequest Story) (err error)
	List(ctx context.Context, status string) (stories []Story, err error)
	CreateScene(ctx context.Context, request CreateSceneRequest) (err error)
	InsertImage(ctx context.Context, request InsertImage) (err error)
	UpdateSceneOrder(ctx context.Context, sceneID string, sceneNumber int64, storyID string) (scene Scene, err error)
}
