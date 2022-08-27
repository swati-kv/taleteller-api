package story

type CreateStoryRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Mood        string `json:"mood,omitempty"`
	Category    string `json:"category,omitempty"`
}

type CreateSceneRequest struct {
	Prompt          string `json:"prompt"`
	Audio           string `json:"audio"`
	BackgroundMusic string `json:"background_music"`
}
