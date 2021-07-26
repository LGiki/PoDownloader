package podcast

// Category is the category of Podcast, if a category has no subcategories, the subcategory will be nil
type Category struct {
	Category    string `json:"category,omitempty"`
	SubCategory string `json:"subCategory,omitempty"`
}
