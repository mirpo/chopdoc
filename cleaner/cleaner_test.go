package cleaner

import (
	"testing"

	"github.com/mirpo/chopdoc/config"
	"github.com/stretchr/testify/assert"
)

func TestCleanText(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		cleanMode config.CleaningMode
		want      string
	}{
		{
			name:      "clean none",
			text:      " test \n\n\ntext\n\n\n\n more ",
			cleanMode: config.CleanNone,
			want:      " test \n\n\ntext\n\n\n\n more ",
		},
		{
			name:      "clean trim",
			text:      "   test   ",
			cleanMode: config.CleanTrim,
			want:      "test",
		},
		{
			name:      "clean normal",
			text:      " test \n\n\ntext\n\n\n\n more ",
			cleanMode: config.CleanNormal,
			want:      "test \ntext\n more",
		},
		{
			name:      "clean aggressive",
			text:      " test \n\n\ntext\n\n\n\n more ",
			cleanMode: config.CleanAggressive,
			want:      "test text more",
		},
		{
			name:      "empty string",
			text:      "",
			cleanMode: config.CleanAggressive,
			want:      "",
		},
		{
			name:      "only whitespace",
			text:      "   \n \t   \n\n   ",
			cleanMode: config.CleanAggressive,
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clean(tt.text, tt.cleanMode)
			assert.Equal(t, tt.want, got)
		})
	}
}
