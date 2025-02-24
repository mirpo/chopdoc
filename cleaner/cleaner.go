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
	case config.CleanNone:
		return chunk
	case config.CleanAggressive:
		chunk = aggressive.ReplaceAllString(chunk, " ")
		fallthrough
	case config.CleanNormal:
		chunk = consecutiveNewlines.ReplaceAllString(chunk, "\n")
		fallthrough
	default:
		return strings.TrimSpace(chunk)
	}
}
