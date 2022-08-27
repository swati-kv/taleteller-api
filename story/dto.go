package story

type CreateStoryRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Mood        string `json:"mood,omitempty"`
	Category    string `json:"category,omitempty"`
}
