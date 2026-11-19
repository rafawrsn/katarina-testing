package utils

import (
	"regexp"
	"strings"
)

var (
	slugNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)
	slugTrimDash = regexp.MustCompile(`^-+|-+$`)
)

// Slugify converts a display name into a URL-safe slug,
// e.g. "Men's Basic Tee" -> "mens-basic-tee".
func Slugify(input string) string {
	s := strings.ToLower(strings.TrimSpace(input))
	s = strings.ReplaceAll(s, "'", "")
	s = slugNonAlnum.ReplaceAllString(s, "-")
	s = slugTrimDash.ReplaceAllString(s, "")
	if s == "" {
		s = "item"
	}
	return s
}
