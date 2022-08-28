package store

import (
	"context"
	"fmt"
)

func (s storyStore) UpdateSceneOrder(ctx context.Context, sceneID string, sceneNumber int64, storyID string) (scene Scene, err error) {
	err = s.db.GetContext(ctx, &scene, updateSceneOrder, sceneID, sceneNumber, storyID)
	if err != nil {
		fmt.Println("error - 1")
		return
	}
	fmt.Println("val - 1 ", scene)
	err = s.db.GetContext(ctx, &scene.SelectedImagePath, getImagePath, scene.SelectedImage)
	if err != nil {
		fmt.Println("error -2")
		return
	}
	err = s.db.GetContext(ctx, &scene.GeneratedAudioPath, getGeneratedAudio, scene.GeneratedAudioID)
	if err != nil {
		fmt.Println("error - 3")
	}
	return
}
