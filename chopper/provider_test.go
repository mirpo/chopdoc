package chopper

import (
	"bufio"
	"strings"
	"testing"

	"github.com/mirpo/chopdoc/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChopper(t *testing.T) {
	tests := []struct {
		name           string
		method         config.ChunkMethod
		cfg            *config.Config
		expectType     string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:       "char chopper",
			method:     config.Char,
			cfg:        &config.Config{ChunkSize: 100},
			expectType: "*chopper.CharChopper",
		},
		{
			name:       "word chopper",
			method:     config.Word,
			cfg:        &config.Config{ChunkSize: 10},
			expectType: "*chopper.WordChopper",
		},
		{
			name:       "sentence chopper",
			method:     config.Sentence,
			cfg:        &config.Config{ChunkSize: 5},
			expectType: "*chopper.SentenceChopper",
		},
		{
			name:       "recursive chopper",
			method:     config.Recursive,
			cfg:        &config.Config{ChunkSize: 100},
			expectType: "*chopper.RecursiveChopper",
		},
		{
			name:       "markdown chopper",
			method:     config.Markdown,
			cfg:        &config.Config{ChunkSize: 100, MarkdownLevels: []int{1, 2, 3}},
			expectType: "*chopper.MarkdownChopper",
		},
		{
			name:           "invalid method",
			method:         config.ChunkMethod("invalid"),
			cfg:            &config.Config{ChunkSize: 100},
			expectError:    true,
			expectedErrMsg: "unsupported chunkMethod: invalid",
		},
		{
			name:           "empty method",
			method:         config.ChunkMethod(""),
			cfg:            &config.Config{ChunkSize: 100},
			expectError:    true,
			expectedErrMsg: "unsupported chunkMethod: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader("test input")
			var output strings.Builder
			rw := bufio.NewReadWriter(bufio.NewReader(input), bufio.NewWriter(&output))

			chopper, err := NewChopper(tt.method, tt.cfg, rw)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, chopper)
				if tt.expectedErrMsg != "" {
					assert.EqualError(t, err, tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, chopper)
				actualType := assert.IsType(t, chopper, chopper)
				if actualType {
					switch tt.method {
					case config.Char:
						assert.IsType(t, &CharChopper{}, chopper)
					case config.Word:
						assert.IsType(t, &WordChopper{}, chopper)
					case config.Sentence:
						assert.IsType(t, &SentenceChopper{}, chopper)
					case config.Recursive:
						assert.IsType(t, &RecursiveChopper{}, chopper)
					case config.Markdown:
						assert.IsType(t, &MarkdownChopper{}, chopper)
					}
				}
			}
		})
	}
}

func TestChopperEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		method config.ChunkMethod
		input  string
		cfg    *config.Config
	}{
		{
			name:   "empty input - char",
			method: config.Char,
			input:  "",
			cfg:    &config.Config{ChunkSize: 10, Overlap: 0},
		},
		{
			name:   "empty input - word",
			method: config.Word,
			input:  "",
			cfg:    &config.Config{ChunkSize: 5, Overlap: 0},
		},
		{
			name:   "empty input - sentence",
			method: config.Sentence,
			input:  "",
			cfg:    &config.Config{ChunkSize: 3, Overlap: 0},
		},
		{
			name:   "whitespace only - char",
			method: config.Char,
			input:  "   \n\t  ",
			cfg:    &config.Config{ChunkSize: 5, Overlap: 0, CleaningMode: config.CleanNone},
		},
		{
			name:   "unicode characters - char",
			method: config.Char,
			input:  "Hello 世界 🌍",
			cfg:    &config.Config{ChunkSize: 5, Overlap: 0},
		},
		{
			name:   "very long line - word",
			method: config.Word,
			input:  strings.Repeat("word ", 1000),
			cfg:    &config.Config{ChunkSize: 10, Overlap: 2},
		},
		{
			name:   "single character - char",
			method: config.Char,
			input:  "a",
			cfg:    &config.Config{ChunkSize: 10, Overlap: 0},
		},
		{
			name:   "exact chunk size - word",
			method: config.Word,
			input:  "one two three four five",
			cfg:    &config.Config{ChunkSize: 5, Overlap: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			var output strings.Builder
			rw := bufio.NewReadWriter(bufio.NewReader(input), bufio.NewWriter(&output))

			chopper, err := NewChopper(tt.method, tt.cfg, rw)
			require.NoError(t, err)
			require.NotNil(t, chopper)

			err = chopper.Chop()
			assert.NoError(t, err)

			err = rw.Flush()
			assert.NoError(t, err)

			outputStr := output.String()
			if outputStr != "" {
				lines := strings.Split(strings.TrimSpace(outputStr), "\n")
				for i, line := range lines {
					if strings.TrimSpace(line) != "" {
						assert.True(t, strings.HasPrefix(line, "{"),
							"Line %d should start with '{': %s", i+1, line)
						assert.True(t, strings.HasSuffix(line, "}"),
							"Line %d should end with '}': %s", i+1, line)
						assert.Contains(t, line, "\"chunk\":",
							"Line %d should contain chunk field: %s", i+1, line)
					}
				}
			}
		})
	}
}
