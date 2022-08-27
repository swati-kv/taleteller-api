package store

import (
	"context"
)

func (s storyStore) UpdateSceneOrder(ctx context.Context, sceneID string, sceneNumber int64, storyID string) (scene Scene, err error) {
	err = s.db.GetContext(ctx, &scene, updateSceneOrder, sceneID, sceneNumber, storyID)
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &scene.SelectedImagePath, getImagePath, scene.SelectedImage)
	if err != nil {
		return
	}
	err = s.db.GetContext(ctx, &scene.GeneratedAudioPath, getGeneratedAudio, scene.GeneratedAudioID)
	return
}
