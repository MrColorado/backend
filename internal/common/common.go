package common

import "strings"

func HarmonizeTitle(title string) string {
	return strings.ToLower(strings.TrimSpace(title))
}
