package story

type CreateStoryRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Mood        string `json:"mood,omitempty"`
	Category    string `json:"category,omitempty"`
	CustomerID  string `json:"customer_id"`
}
type Image struct {
	SelectedImage string `json:"selected_image"`
}

type CreateSceneRequest struct {
	Prompt          string `json:"prompt"`
	Audio           string `json:"audio"`
	BackgroundMusic string `json:"background_music"`
	ImageCount      int64  `json:"image_count"`
	SceneNumber     int64  `json:"scene_number"`
}

type CreateSceneResponse struct {
	Status  string `json:"status"`
	SceneID string `json:"scene_id"`
}

type PyImageRequest struct {
	Prompt string `json:"prompt"`
	Count  int64  `json:"num"`
}

type PyImageResponse struct {
	Data struct {
		GeneratedImage       []string `json:"generatedImgs"`
		GeneratedImageFormat string   `json:"generatedImgsFormat"`
	} `json:"data"`
	Error string `json:"error"`
}

type PyAudioRequest struct {
	Prompt   string `json:"prompt"`
	Language string `json:"lang"`
}

type PyAudioResponse struct {
	Data  string `json:"data"`
	Error string `json:"error"`
}

type ImageDetails struct {
	ImageID   string `json:"image_id"`
	ImagePath string `json:"image_path"`
}

type GetSceneResponse struct {
	Status string         `json:"status"`
	Images []ImageDetails `json:"images"`
}
