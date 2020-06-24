package util

func SubString(s string, start, end int) string {
	souce := []rune(s)
	n := len(souce)
	if start < 0 || end > n || start > end {
		return ""
	}
	return string(souce[start:end])
}
