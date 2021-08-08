package podcast

import "encoding/json"

// Category is the category of Podcast, if a category has no subcategories, the SubCategory will be nil
type Category struct {
	Category    string `json:"category,omitempty"`
	SubCategory string `json:"subCategory,omitempty"`
}

// GetJSON returns a Category instance in JSON format
func (c *Category) GetJSON() (string, error) {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
