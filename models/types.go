//types.go

package models

// CatImage is a struct representing a cat image response from The Cat API
type CatImage struct {
	ID     string  `json:"id"`
	URL    string  `json:"url"`
	Breeds []Breed `json:"breeds"`
}

// Breed is a struct for cat breed information
type Breed struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Origin       string `json:"origin"`
	Temperament  string `json:"temperament"`
	WikipediaURL string `json:"wikipedia_url"`
}

// Vote is used for submitting or retrieving votes from The Cat API
type Vote struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id"`
	Value   int    `json:"value"` // 1 for up, -1 for down
}

// FavoritePayload represents the request body to add a favorite
type FavoritePayload struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id"`
}
