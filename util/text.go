package util

// GetFirstNCharacters
// Get first n characters of string with Chinese support,
// if giving n is greater then string length will return string directly
func GetFirstNCharacters(text string, n int) string {
	textRune := []rune(text)
	if len(textRune) >= n {
		return string(textRune[:n])
	}
	return text
}

// FillTextToLength
// Fill the text to giving length by append blanks to the string,
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
