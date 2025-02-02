package runner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mirpo/chopdoc/chopper"
	"github.com/mirpo/chopdoc/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChar(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		input      string
		chunkSize  int
		overlap    int
		cleanMode  config.CleaningMode
		wantChunks []string
		wantErr    bool
	}{
		{
			name:       "basic chunking",
			input:      "hello world test",
			chunkSize:  5,
			overlap:    0,
			wantChunks: []string{"hello", " worl", "d tes", "t"},
		},
		{
			name:       "with overlap",
			input:      "hello world",
			chunkSize:  6,
			overlap:    2,
			wantChunks: []string{"hello ", "o worl", "rld"},
		},
		{
			name:      "empty file",
			input:     "",
			chunkSize: 5,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inPath := filepath.Join(tmpDir, "input.txt")
			err := os.WriteFile(inPath, []byte(tt.input), 0o644)
			require.NoError(t, err)

			outPath := filepath.Join(tmpDir, "output.jsonl")

			cfg := &config.Config{
				InputFile:    inPath,
				OutputFile:   outPath,
				ChunkSize:    tt.chunkSize,
				Overlap:      tt.overlap,
				CleaningMode: tt.cleanMode,
				Method:       config.Char,
			}

			r := NewRunner(cfg)
			err = r.Run()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			f, err := os.Open(outPath)
			require.NoError(t, err)
			defer f.Close()

			var chunks []chopper.Chunk
			dec := json.NewDecoder(f)
			for dec.More() {
				var chunk chopper.Chunk
				require.NoError(t, dec.Decode(&chunk))
				chunks = append(chunks, chunk)
			}

			assert.Equal(t, len(tt.wantChunks), len(chunks))
			for i, want := range tt.wantChunks {
				assert.Equal(t, want, chunks[i].Text)
			}
		})
	}
}

func TestErrorCases(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.txt")
	err := os.WriteFile(inputPath, []byte("test content"), 0o644)
	require.NoError(t, err)

	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr string
	}{
		{
			name: "input file not found",
			cfg: &config.Config{
				InputFile:  "nonexistent.txt",
				OutputFile: "out.jsonl",
				ChunkSize:  10,
			},
			wantErr: "failed to open input file",
		},
		{
			name: "invalid output path",
			cfg: &config.Config{
				InputFile:  inputPath,
				OutputFile: "/invalid/path/out.jsonl",
				ChunkSize:  10,
			},
			wantErr: "failed to create output file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRunner(tt.cfg)
			err := r.Run()
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}
