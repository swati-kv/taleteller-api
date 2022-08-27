package store

const (
	updateScene = `UPDATE scene SET scene_number=$2 WHERE id=$1 AND story_id=$3 RETURNING generated_audio_id,background_audio_path,selected_image`

	getGeneratedAudio = `SELECT path from generated_audio WHERE id=$1`

	getImagePath = `SELECT path from image WHERE id=$1`
)
