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
