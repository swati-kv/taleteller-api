package store

import (
	"context"
	"github.com/jmoiron/sqlx"
	"taleteller/logger"
	"time"
)

type storyStore struct {
	db *sqlx.DB
}

func (s storyStore) InsertImage(ctx context.Context, request InsertImageRequest) (err error) {
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

func (s storyStore) List(ctx context.Context, status string) (stories []Story, err error) {
	err = s.db.SelectContext(ctx, &stories, getStories, status)
	return
}

func (s storyStore) Create(ctx context.Context, c Story) (err error) {
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
func (s storyStore) GetStoryByID(ctx context.Context, storyID string) (storyResponse Story, err error) {
	err = s.db.GetContext(ctx, &storyResponse, getStoryByID, storyID)
	err = s.db.SelectContext(ctx, &storyResponse.SceneDetails, getSceneByID, storyID)
	return
}
func (s storyStore) UpdateScene(ctx context.Context, storyID string, sceneID string, selectedImage string) (scene Scene, err error) {
	err = s.db.GetContext(ctx, &scene, updateScene, selectedImage, storyID, sceneID)
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

func (s *storyStore) InsertAudio(ctx context.Context, request InsertAudioRequest) (err error) {
	_, err = s.db.ExecContext(ctx, insertAudio,
		request.ID,
		request.AudioPath,
		time.Now(),
		time.Now(),
	)
	return
}

func (s *storyStore) UpdateSceneAudio(ctx context.Context, id string, sceneID string) (err error) {
	_, err = s.db.ExecContext(ctx, updateAudioInScene,
		id,
		sceneID,
		time.Now(),
	)
	return
}

func (s *storyStore) GetSceneStatus(ctx context.Context, sceneID string) (status string, err error) {
	err = s.db.SelectContext(ctx, &status, getSceneStatusByID, sceneID)
	return

}

func (s *storyStore) UpdateSceneStatus(ctx context.Context, media string, sceneID string, incomingStatus string) (err error) {
	currentStatus, err := s.GetSceneStatus(ctx, sceneID)
	if err != nil {
		logger.Errorw(ctx, "error getting scene status")
		return
	}
	var realStatus string
	switch media {
	case "image":
		if currentStatus == "audio_done" {
			realStatus = "completed"
		} else {
			realStatus = incomingStatus
		}
		_, err = s.db.ExecContext(ctx, updateMediaStatusInScene,
			realStatus,
			sceneID,
			time.Now(),
		)
		if err != nil {
			return
		}
	case "audio":
		if currentStatus == "image_done" {
			realStatus = "completed"
		} else {
			realStatus = incomingStatus
		}
		_, err = s.db.ExecContext(ctx, updateMediaStatusInScene,
			realStatus,
			sceneID,
			time.Now(),
		)
		if err != nil {
			return
		}
	}

	return
}

func (s *storyStore) GetSceneByID(ctx context.Context, sceneID string, storyID string) (response []GetSceneByIDResponse, err error) {
	err = s.db.SelectContext(ctx, &response, getSceneDetailsByID, sceneID, storyID)
	return
}

func NewStoryStore(db *sqlx.DB) *storyStore {
	return &storyStore{
		db: db,
	}
}
