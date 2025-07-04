package runner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mirpo/chopdoc/chopper"
	"github.com/mirpo/chopdoc/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestValidatePathIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.txt")
	err := os.WriteFile(inputPath, []byte("test content"), 0o644)
	require.NoError(t, err)

	tests := []struct {
		name       string
		inputFile  string
		outputFile string
		wantErr    string
	}{
		{
			name:       "valid paths",
			inputFile:  inputPath,
			outputFile: filepath.Join(tmpDir, "output.jsonl"),
			wantErr:    "",
		},
		{
			name:       "path traversal in input",
			inputFile:  "../../../etc/passwd",
			outputFile: filepath.Join(tmpDir, "output.jsonl"),
			wantErr:    "invalid input file path: path traversal detected",
		},
		{
			name:       "path traversal in output",
			inputFile:  inputPath,
			outputFile: "../../../tmp/malicious.jsonl",
			wantErr:    "invalid output file path: path traversal detected",
		},
		{
			name:       "nested path traversal",
			inputFile:  "some/../../../etc/passwd",
			outputFile: filepath.Join(tmpDir, "output.jsonl"),
			wantErr:    "invalid input file path: path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				InputFile:  tt.inputFile,
				OutputFile: tt.outputFile,
				ChunkSize:  10,
				Method:     config.Char,
			}

			r := NewRunner(cfg)
			err := r.Run()

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestRunnerWithInvalidChopper(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.txt")
	err := os.WriteFile(inputPath, []byte("test content"), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		InputFile:  inputPath,
		OutputFile: filepath.Join(tmpDir, "output.jsonl"),
		ChunkSize:  10,
		Method:     config.ChunkMethod("invalid"), // Invalid method
	}

	r := NewRunner(cfg)
	err = r.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create chopper")
}

func TestChar(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		input      string
		chunkSize  int
		overlap    int
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
				Method:       config.Char,
				CleaningMode: config.CleanNone,
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

func TestWord(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		input      string
		chunkSize  int
		overlap    int
		wantChunks []string
		wantErr    bool
	}{
		{
			name:      "basic chunking",
			input:     "one two three four five six seven eight nine ten.",
			chunkSize: 3,
			overlap:   0,
			wantChunks: []string{
				"one two three",
				"four five six",
				"seven eight nine",
				"ten.",
			},
		},
		{
			name:      "with overlap",
			input:     "one two three four five six seven eight nine ten.",
			chunkSize: 3,
			overlap:   1,
			wantChunks: []string{
				"one two three",
				"three four five",
				"five six seven",
				"seven eight nine",
				"nine ten.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inPath := filepath.Join(tmpDir, "input.txt")
			err := os.WriteFile(inPath, []byte(tt.input), 0o644)
			require.NoError(t, err)

			outPath := filepath.Join(tmpDir, "output.jsonl")

			cfg := &config.Config{
				InputFile:  inPath,
				OutputFile: outPath,
				ChunkSize:  tt.chunkSize,
				Overlap:    tt.overlap,
				Method:     config.Word,
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

func TestSentence(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		input      string
		chunkSize  int
		overlap    int
		wantChunks []string
		wantErr    bool
	}{
		{
			name:      "basic chunking one",
			input:     "basic chunking one.   chunking two? chunking three!.",
			chunkSize: 1,
			overlap:   0,
			wantChunks: []string{
				"basic chunking one.",
				"chunking two?",
				"chunking three!.",
			},
		},
		{
			name:      "basic chunking one 2",
			input:     "basic chunking one.   chunking two? chunking three!.",
			chunkSize: 2,
			overlap:   0,
			wantChunks: []string{
				"basic chunking one. chunking two?",
				"chunking three!.",
			},
		},
		{
			name:      "with overlap",
			input:     "basic chunking one.   chunking two? chunking three!.",
			chunkSize: 2,
			overlap:   1,
			wantChunks: []string{
				"basic chunking one. chunking two?",
				"chunking two? chunking three!.",
				"chunking three!.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inPath := filepath.Join(tmpDir, "input.txt")
			err := os.WriteFile(inPath, []byte(tt.input), 0o644)
			require.NoError(t, err)

			outPath := filepath.Join(tmpDir, "output.jsonl")

			cfg := &config.Config{
				InputFile:  inPath,
				OutputFile: outPath,
				ChunkSize:  tt.chunkSize,
				Overlap:    tt.overlap,
				Method:     config.Sentence,
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

func TestRecursive(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		input      string
		chunkSize  int
		overlap    int
		wantChunks []string
		wantErr    bool
	}{
		{
			name:      "basic chunking one",
			input:     "basic chunking one.   chunking two? chunking three!.",
			chunkSize: 20,
			overlap:   0,
			wantChunks: []string{
				"basic chunking one.",
				"chunking two?",
				"chunking three!.",
			},
		},
		{
			name: "basic chunking two",
			input: `basic chunking one.
			
			chunking two?
			chunking three!.
			
			
			chunking four!.`,
			chunkSize: 20,
			overlap:   0,
			wantChunks: []string{
				"basic chunking one.",
				"chunking two?",
				"chunking three!.",
				"chunking four!.",
			},
		},
		{
			name: "basic chunking three",
			input: `basic chunking one.
			
			chunking two?
			chunking three!.
			
			
			chunking four!.`,
			chunkSize: 50,
			overlap:   0,
			wantChunks: []string{
				"basic chunking one.\n\t\t\t\n\t\t\tchunking two?",
				"chunking three!.",
				"chunking four!.",
			},
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
				Method:       config.Recursive,
				CleaningMode: config.CleanTrim,
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

func TestMarkdownChopper(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		levels     []int
		strip      bool
		addMeta    bool
		wantChunks []chopper.Chunk
	}{
		{
			name:    "Basic Markdown Chunking",
			input:   "# Header 1\nContent under header 1\n\n## Header 2\nContent under header 2\n",
			levels:  []int{1, 2},
			strip:   false,
			addMeta: true,
			wantChunks: []chopper.Chunk{
				{Text: "# Header 1\nContent under header 1\n", Metadata: map[string]string{"Header 1": "Header 1"}},
				{Text: "## Header 2\nContent under header 2\n", Metadata: map[string]string{"Header 1": "Header 1", "Header 2": "Header 2"}},
			},
		},
		{
			name:    "Stripping Headers",
			input:   "# Header 1\nContent under header 1\n\n## Header 2\nContent under header 2\n",
			levels:  []int{1, 2},
			strip:   true,
			addMeta: true,
			wantChunks: []chopper.Chunk{
				{Text: "Content under header 1\n", Metadata: map[string]string{"Header 1": "Header 1"}},
				{Text: "Content under header 2\n", Metadata: map[string]string{"Header 1": "Header 1", "Header 2": "Header 2"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			rw := bufio.NewReadWriter(bufio.NewReader(strings.NewReader(tt.input)), bufio.NewWriter(&output))

			cfg := &config.Config{
				MarkdownLevels: tt.levels,
				StripHeaders:   tt.strip,
				AddMetadata:    tt.addMeta,
			}

			markdownChopper := chopper.NewMarkdownChopper(cfg, rw)
			err := markdownChopper.Chop()

			require.NoError(t, err)
			require.NoError(t, rw.Flush())

			var gotChunks []chopper.Chunk
			decoder := json.NewDecoder(&output)

			for decoder.More() {
				var chunk chopper.Chunk
				require.NoError(t, decoder.Decode(&chunk))
				gotChunks = append(gotChunks, chunk)
			}

			assert.Equal(t, tt.wantChunks, gotChunks)
		})
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	fmt.Print("")
	return buf.String()
}

func TestPiped(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		chunkSize           int
		overlap             int
		expectedPipedOutput string
		wantErr             bool
	}{
		{
			name:                "basic chunking one",
			input:               "basic chunking one.   chunking two? chunking three!.",
			chunkSize:           1,
			overlap:             0,
			expectedPipedOutput: "{\"chunk\":\"basic chunking one.\"}\n{\"chunk\":\"chunking two?\"}\n{\"chunk\":\"chunking three!.\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			r, w, _ := os.Pipe()
			os.Stdin = r

			cfg := &config.Config{
				ChunkSize: tt.chunkSize,
				Overlap:   tt.overlap,
				Method:    config.Sentence,
				Piped:     true,
			}
			testInput := tt.input
			go func() {
				_, _ = w.Write([]byte(testInput))
				w.Close()
			}()
			var err error

			runner := NewRunner(cfg)
			output := captureOutput(func() {
				err = runner.Run()
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
			})

			require.NoError(t, err)

			assert.Equal(t, tt.expectedPipedOutput, output)
		})
	}
}
