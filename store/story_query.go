package store

const (
	createStory = `INSERT INTO public.story
(id, name, mood, category, description, customer_id, status, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6, $7, $8,$9);
`
	getStoryByID = `SELECT * FROM public.story WHERE id = $1`

	getSceneByID = `SELECT sc.id, sc.story_id, sc.generated_audio_id, sc.selected_image, g.path, sc.background_audio_path, sc.status, sc.scene_number, sc.created_at, sc.updated_at from public.scene sc INNER JOIN public.generated_audio g ON sc.generated_audio_id = g.id WHERE story_id = $1`

	getStories = `SELECT id, "name", mood, category, description, customer_id, status, created_at, updated_at
FROM story WHERE status = $1`

	updateScene = `UPDATE public.scene SET selected_image = $1 where story_id = $2 and id = $3 RETURNING *`
)
