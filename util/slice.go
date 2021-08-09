package util

// IsStringSliceContainText returns true if specified string slice s contains specified text,
// otherwise it returns false
func IsStringSliceContainText(s []string, text string) bool {
	for _, v := range s {
		if v == text {
			return true
		}
	}
	return false
}
