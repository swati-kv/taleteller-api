package store

const (
	createStory = `INSERT INTO public.story
(id, name, mood, category, description, customer_id, status, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6, $7, $8,$9);`

	createScene = `INSERT INTO public.scene
(id, story_id, status, scene_number, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6);
`

	insertImage = `INSERT INTO public.image
(id, "path", scene_id, is_deleted, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6);
`
	getStoryByID = `SELECT * FROM public.story WHERE id = $1`

	getSceneByID = `SELECT sc.id, sc.story_id, sc.generated_audio_id, sc.selected_image, g.path, sc.background_audio_path, sc.status, sc.scene_number, sc.created_at, sc.updated_at from public.scene sc INNER JOIN public.generated_audio g ON sc.generated_audio_id = g.id WHERE story_id = $1`

	getStories = `SELECT id, "name", mood, category, description, customer_id, status, created_at, updated_at
FROM story WHERE status = $1`

	updateScene = `UPDATE public.scene SET selected_image = $1 where story_id = $2 and id = $3 RETURNING *`
	insertAudio = `INSERT INTO public.generated_audio
(id, "path", created_at, updated_at)
VALUES($1, $2, $3, $4);
`
	updateAudioInScene = `UPDATE public.scene
SET generated_audio_id=$1, updated_at=$3
WHERE id=$2;
`
	updateMediaStatusInScene = `UPDATE public.scene
SET status=$1, updated_at=$3
WHERE id=$2;
`
	getSceneStatusByID = `SELECT status from public.scene where id = $1`

	getSceneDetailsByID = `select i.id, path, status from image i inner join scene s  on i.scene_id = s.id where s.story_id = $2 and s.id = $1;`
)
