package runner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mirpo/chopdoc/cleaner"
	"github.com/mirpo/chopdoc/config"
)

type Chunk struct {
	Text string `json:"chunk"`
}

type Runner struct {
	cfg *config.Config
}

func NewRunner(cfg *config.Config) *Runner {
	return &Runner{
		cfg: cfg,
	}
}

func (r *Runner) cleanText(text string) string {
	return cleaner.Clean(text, r.cfg.CleaningMode)
}

func (r *Runner) writeChunk(encoder *json.Encoder, chunk string) error {
	chunkText := r.cleanText(chunk)

	if len(strings.TrimSpace(chunkText)) == 0 {
		return nil
	}

	if err := encoder.Encode(Chunk{Text: chunkText}); err != nil {
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	return nil
}

func (r *Runner) Run() error {
	input, err := os.Open(r.cfg.InputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	output, err := os.Create(r.cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()

	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)
	rw := bufio.NewReadWriter(reader, writer)

	encoder := json.NewEncoder(rw.Writer)

	scanner := bufio.NewScanner(rw.Reader)
	scanner.Split(bufio.ScanBytes)

	chunk := ""
	step := r.cfg.ChunkSize - r.cfg.Overlap

	for scanner.Scan() {
		chunk += scanner.Text()

		if len(chunk) >= r.cfg.ChunkSize {
			err := r.writeChunk(encoder, chunk)
			if err != nil {
				return err
			}

			chunk = chunk[step:]
		}
	}

	if len(chunk) > 0 {
		err := r.writeChunk(encoder, chunk)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input file:", err)
	}

	rw.Writer.Flush()

	return nil
}
