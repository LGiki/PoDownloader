package util

// GetFirstNCharacters returns first n characters of string with Chinese support,
// if giving n is greater than string length will return string directly
func GetFirstNCharacters(text string, n int) string {
	textRune := []rune(text)
	if len(textRune) >= n {
		return string(textRune[:n])
	}
	return text
}

// FillTextToLength fills the text to giving length by append blanks to the string,
// if the string length the greater than giving length will return string directly
func FillTextToLength(text string, length int) string {
	if len(text) < length {
		appendCount := length - len(text)
		for i := 0; i < appendCount; i++ {
			text += " "
		}
	}
	return text
}

// RemoveDuplicateItemsInStringSlice returns the string slice after removing duplicates and duplicated string items
func RemoveDuplicateItemsInStringSlice(texts []string) (result []string, duplicatedItems []string) {
	duplicated := make(map[string]bool)
	for _, text := range texts {
		_, isExist := duplicated[text]
		if isExist {
			duplicatedItems = append(duplicatedItems, text)
		} else {
			duplicated[text] = true
			result = append(result, text)
		}
	}
	return
}
