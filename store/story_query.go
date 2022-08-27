package store

const (
	createStory = `INSERT INTO public.story
(id, name, mood, category, description, customer_id, status, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6, $7, $8,$9);
`
	createScene = `INSERT INTO public.scene
(id, story_id, status, scene_number, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6);
`

	insertImage = `INSERT INTO public.image
(id, "path", scene_id, is_deleted, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6);
`
)
