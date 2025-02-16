package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, 1000, cfg.ChunkSize)
	assert.Equal(t, CleanNormal, cfg.CleaningMode)
	assert.Equal(t, 0, cfg.Overlap)
	assert.Equal(t, false, cfg.Piped)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
