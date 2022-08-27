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
	GeneratedAudioPath  string     `json:"generated-audio-path,omitempty" db:"path,omitempty"`
	BackgroundAudioPath string     `db:"background_audio_path,omitempty" json:"background-audio-path,omitempty"`
	SelectedImage       string     `db:"selected_image,omitempty" json:"selected_image"`
	Status              string     `db:"status,omitempty" json:"status,omitempty"`
	SceneNumber         int        `db:"scene_number,omitempty" json:"scene-number,omitempty"`
	SelectedImagePath   string     `db:"path" json:"selected-image-path,omitempty"`
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
	SceneID         string
	Status          string
	StoryID         string
	SceneNumber     int64
	BackgroundMusic string
}

type InsertImageRequest struct {
	ID        string
	ImagePath string
	SceneID   string
}

type InsertAudioRequest struct {
	ID        string
	AudioPath string
	SceneID   string
}

type GetSceneByIDResponse struct {
	ImageID   string `db:"id"`
	ImagePath string `db:"path"`
	Status    string `db:"status"`
}

type StoryStorer interface {
	GetStoryByID(ctx context.Context, storyID string) (storyDetails Story, err error)
	Create(ctx context.Context, createRequest Story) (err error)
	List(ctx context.Context, status string) (stories []Story, err error)
	UpdateScene(ctx context.Context, storyID string, sceneID string, selectedImage string) (sceneDetails Scene, err error)
	CreateScene(ctx context.Context, request CreateSceneRequest) (err error)
	UpdateSceneOrder(ctx context.Context, sceneID string, sceneNumber int64, storyID string) (scene Scene, err error)
	InsertImage(ctx context.Context, request InsertImageRequest) (err error)
	InsertAudio(ctx context.Context, request InsertAudioRequest) (err error)
	UpdateSceneAudio(background context.Context, id string, sceneID string) (err error)
	UpdateSceneStatus(background context.Context, media string, sceneID string, status string) (err error)
	GetSceneByID(ctx context.Context, id string, id2 string) (response []GetSceneByIDResponse, err error)
}
