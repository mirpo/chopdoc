package cleaner

import (
	"regexp"
	"strings"

	"github.com/mirpo/chopdoc/config"
)

var (
	aggressive          = regexp.MustCompile(`[\p{Z}\p{C}\s]+`)
	consecutiveNewlines = regexp.MustCompile(`\n\s*\n+`)
)

func Clean(chunk string, cleaningMode config.CleaningMode) string {
	switch cleaningMode {
	case config.CleanAggressive:
		chunk = aggressive.ReplaceAllString(chunk, " ")
		fallthrough
	case config.CleanNormal:
		chunk = consecutiveNewlines.ReplaceAllString(chunk, "\n")
		fallthrough
	case config.CleanTrim:
		chunk = strings.TrimSpace(chunk)
		fallthrough
	default:
		return chunk
	}
}
