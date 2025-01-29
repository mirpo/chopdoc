// runner/runner_test.go
package runner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mirpo/chopdoc/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRunner(&config.Config{CleaningMode: tt.cleanMode})
			got := r.cleanText(tt.text)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRun(t *testing.T) {
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
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create input file
			inPath := filepath.Join(tmpDir, "input.txt")
			err := os.WriteFile(inPath, []byte(tt.input), 0o644)
			require.NoError(t, err)

			// Setup output file
			outPath := filepath.Join(tmpDir, "output.jsonl")

			cfg := &config.Config{
				InputFile:    inPath,
				OutputFile:   outPath,
				ChunkSize:    tt.chunkSize,
				Overlap:      tt.overlap,
				CleaningMode: tt.cleanMode,
			}

			r := NewRunner(cfg)
			err = r.Run()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Read and verify chunks
			f, err := os.Open(outPath)
			require.NoError(t, err)
			defer f.Close()

			var chunks []Chunk
			dec := json.NewDecoder(f)
			for dec.More() {
				var chunk Chunk
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
	// Create test input file
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
