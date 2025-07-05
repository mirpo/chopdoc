package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, 1000, cfg.ChunkSize)
	assert.Equal(t, CleanNone, cfg.CleaningMode)
	assert.Equal(t, 0, cfg.Overlap)
	assert.Equal(t, false, cfg.Piped)
	assert.Equal(t, "1-6", cfg.MarkdownHeader)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, cfg.MarkdownLevels)
	assert.Equal(t, false, cfg.StripHeaders)
	assert.Equal(t, false, cfg.AddMetadata)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr string
	}{
		{
			name: "valid config",
			cfg: Config{
				InputFile:  "input.txt",
				OutputFile: "output.jsonl",
				Method:     Char,
				ChunkSize:  1000,
				Overlap:    100,
			},
		},
		{
			name: "missing input",
			cfg: Config{
				OutputFile: "output.jsonl",
				ChunkSize:  1000,
			},
			wantErr: "input file is required",
		},
		{
			name: "input can be empty when piped",
			cfg: Config{
				Piped:      true,
				OutputFile: "",
				ChunkSize:  1000,
				Method:     Char,
			},
		},
		{
			name: "empty output is allowed",
			cfg: Config{
				InputFile: "input.txt",
				ChunkSize: 1000,
				Method:    Char,
			},
		},
		{
			name: "non empty output must include jsonl",
			cfg: Config{
				Piped:      true,
				OutputFile: "output.json",
				ChunkSize:  1000,
				Method:     Char,
			},
			wantErr: "output file must have .jsonl extension",
		},
		{
			name: "invalid chunk size",
			cfg: Config{
				InputFile:  "input.txt",
				OutputFile: "output.jsonl",
				ChunkSize:  0,
			},
			wantErr: "chunk size must be greater than 0",
		},
		{
			name: "overlap too large",
			cfg: Config{
				InputFile:  "input.txt",
				OutputFile: "output.jsonl",
				ChunkSize:  100,
				Overlap:    200,
			},
			wantErr: "overlap must be less than chunk size",
		},
		{
			name: "wrong output extension",
			cfg: Config{
				InputFile:  "input.txt",
				OutputFile: "output.txt",
				ChunkSize:  1000,
			},
			wantErr: "output file must have .jsonl extension",
		},
		{
			name: "valid markdown config",
			cfg: Config{
				InputFile:      "input.md",
				OutputFile:     "output.jsonl",
				Method:         Markdown,
				ChunkSize:      1000,
				MarkdownHeader: "1-6",
			},
		},
		{
			name: "valid markdown config with header range 2-4",
			cfg: Config{
				InputFile:      "input.md",
				OutputFile:     "output.jsonl",
				Method:         Markdown,
				ChunkSize:      1000,
				MarkdownHeader: "2-4",
			},
		},
		{
			name: "invalid markdown header format",
			cfg: Config{
				InputFile:      "input.md",
				OutputFile:     "output.jsonl",
				Method:         Markdown,
				ChunkSize:      1000,
				MarkdownHeader: "invalid",
			},
			wantErr: "invalid markdown header format: invalid, expected format like '1-6'",
		},
		{
			name: "invalid markdown header format with comma",
			cfg: Config{
				InputFile:      "input.md",
				OutputFile:     "output.jsonl",
				Method:         Markdown,
				ChunkSize:      1000,
				MarkdownHeader: "1,3,5",
			},
			wantErr: "invalid markdown header format: 1,3,5, expected format like '1-6'",
		},
		{
			name: "invalid markdown header with start > end",
			cfg: Config{
				InputFile:      "input.md",
				OutputFile:     "output.jsonl",
				Method:         Markdown,
				ChunkSize:      1000,
				MarkdownHeader: "5-2",
			},
			wantErr: "start level (5) must be less than or equal to end level (2)",
		},
		{
			name: "invalid markdown header with level out of range",
			cfg: Config{
				InputFile:      "input.md",
				OutputFile:     "output.jsonl",
				Method:         Markdown,
				ChunkSize:      1000,
				MarkdownHeader: "1-7",
			},
			wantErr: "invalid markdown header format: 1-7, expected format like '1-6'",
		},
		{
			name: "markdown with strip headers enabled",
			cfg: Config{
				InputFile:      "input.md",
				OutputFile:     "output.jsonl",
				Method:         Markdown,
				ChunkSize:      1000,
				MarkdownHeader: "1-6",
				StripHeaders:   true,
			},
		},
		{
			name: "markdown with hide metadata enabled",
			cfg: Config{
				InputFile:      "input.md",
				OutputFile:     "output.jsonl",
				Method:         Markdown,
				ChunkSize:      1000,
				MarkdownHeader: "1-6",
				AddMetadata:    true,
			},
		},
		{
			name: "invalid method type",
			cfg: Config{
				InputFile:  "input.txt",
				OutputFile: "output.jsonl",
				Method:     ChunkMethod("invalid"),
				ChunkSize:  1000,
			},
			wantErr: "invalid chunking method: 'invalid'",
		},
		{
			name: "empty method",
			cfg: Config{
				InputFile:  "input.txt",
				OutputFile: "output.jsonl",
				Method:     ChunkMethod(""),
				ChunkSize:  1000,
			},
			wantErr: "invalid chunking method: ''",
		},
		{
			name: "recursive with overlap shows warning",
			cfg: Config{
				InputFile:  "input.txt",
				OutputFile: "output.jsonl",
				Method:     Recursive,
				ChunkSize:  100,
				Overlap:    50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}

			if tt.cfg.Method == Markdown && tt.wantErr == "" {
				switch tt.cfg.MarkdownHeader {
				case "1-6":
					assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, tt.cfg.MarkdownLevels)
				case "2-4":
					assert.Equal(t, []int{2, 3, 4}, tt.cfg.MarkdownLevels)
				}
			}

			if tt.cfg.Method == Recursive && tt.wantErr == "" && tt.cfg.Overlap != 0 {
				assert.Equal(t, 0, tt.cfg.Overlap, "Recursive method should reset overlap to 0")
			}
		})
	}
}

func TestParseMarkdownHeader(t *testing.T) {
	tests := []struct {
		name           string
		headerStr      string
		expectedLevels []int
		wantErr        bool
		errMsg         string
	}{
		{
			name:           "all headers 1-6",
			headerStr:      "1-6",
			expectedLevels: []int{1, 2, 3, 4, 5, 6},
			wantErr:        false,
		},
		{
			name:           "headers 2-4",
			headerStr:      "2-4",
			expectedLevels: []int{2, 3, 4},
			wantErr:        false,
		},
		{
			name:           "headers 1-1",
			headerStr:      "1-1",
			expectedLevels: []int{1},
			wantErr:        false,
		},
		{
			name:           "headers 6-6",
			headerStr:      "6-6",
			expectedLevels: []int{6},
			wantErr:        false,
		},
		{
			name:      "invalid format",
			headerStr: "invalid",
			wantErr:   true,
			errMsg:    "invalid markdown header format: invalid, expected format like '1-6'",
		},
		{
			name:      "comma-separated not allowed",
			headerStr: "1,3,5",
			wantErr:   true,
			errMsg:    "invalid markdown header format: 1,3,5, expected format like '1-6'",
		},
		{
			name:      "start greater than end",
			headerStr: "5-2",
			wantErr:   true,
			errMsg:    "start level (5) must be less than or equal to end level (2)",
		},
		{
			name:      "out of range",
			headerStr: "0-6",
			wantErr:   true,
			errMsg:    "invalid markdown header format: 0-6, expected format like '1-6'",
		},
		{
			name:      "out of range end",
			headerStr: "1-7",
			wantErr:   true,
			errMsg:    "invalid markdown header format: 1-7, expected format like '1-6'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{MarkdownHeader: tt.headerStr}
			err := cfg.ParseMarkdownHeader()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.EqualError(t, err, tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedLevels, cfg.MarkdownLevels)
			}
		})
	}
}
