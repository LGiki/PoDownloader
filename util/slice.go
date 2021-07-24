package util

func IsStringSliceContainText(s []string, text string) bool {
	for _, v := range s {
		if v == text {
			return true
		}
	}
	return false
}
