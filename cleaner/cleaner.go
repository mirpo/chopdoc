package cleaner

import (
	"regexp"
	"strings"

	"github.com/mirpo/chopdoc/config"
)

var (
	whitespaceCollapse  = regexp.MustCompile(`\s+`)
	consecutiveNewlines = regexp.MustCompile(`\n\s*\n+`)
)

func Clean(chunk string, cleaningMode config.CleaningMode) string {
	switch cleaningMode {
	case config.CleanAggressive:
		chunk = whitespaceCollapse.ReplaceAllString(chunk, " ")
		chunk = consecutiveNewlines.ReplaceAllString(chunk, "\n")
		chunk = strings.TrimSpace(chunk)
	case config.CleanNormal:
		chunk = consecutiveNewlines.ReplaceAllString(chunk, "\n")
		chunk = strings.TrimSpace(chunk)
	case config.CleanTrim:
		chunk = strings.TrimSpace(chunk)
	}
	return chunk
}
