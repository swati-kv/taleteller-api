package store

import (
	"context"
	"taleteller/logger"
)

func (s storyStore) UpdateSceneOrder(ctx context.Context, sceneID string, sceneNumber int64, storyID string) (scene Scene, err error) {
	logger.Info(ctx, "sds ", sceneID, sceneNumber, storyID)
	err = s.db.GetContext(ctx, &scene, updateScene, sceneID, sceneNumber, storyID)
	if err != nil {
		return
	}
	logger.Info(ctx, "here 1 ", scene)
	err = s.db.GetContext(ctx, &scene.SelectedImagePath, getImagePath, scene.SelectedImage)
	if err != nil {
		return
	}
	logger.Info(ctx, "here   2 ", scene)
	err = s.db.GetContext(ctx, &scene.GeneratedAudioPath, getGeneratedAudio, scene.GeneratedAudioID)
	return
}
