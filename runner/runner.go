package runner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

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

var (
	reAggressive          = regexp.MustCompile(`[\p{Z}\p{C}\s]+`)
	reConsecutiveNewlines = regexp.MustCompile(`\n\s*\n+`)
)

func (r *Runner) cleanText(text string) string {
	switch r.cfg.CleaningMode {
	case config.CleanAggressive:
		text = reAggressive.ReplaceAllString(text, " ") // Handle spaces and control chars
		fallthrough
	case config.CleanNormal:
		text = reConsecutiveNewlines.ReplaceAllString(text, "\n") // Remove extra newlines
		return strings.TrimSpace(text)
	default:
		return text
	}
}

func (r *Runner) Run() error {
	// Open input file
	input, err := os.Open(r.cfg.InputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	// Create output file
	output, err := os.Create(r.cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()

	encoder := json.NewEncoder(output)
	reader := bufio.NewReader(input)

	var data []rune
	for {
		char, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading input file: %w", err)
		}
		data = append(data, char)
	}

	if len(data) == 0 {
		return fmt.Errorf("input file is empty")
	}

	// Process chunks with overlap
	chunkSize := r.cfg.ChunkSize
	overlap := r.cfg.Overlap
	step := chunkSize - overlap

	for i := 0; i < len(data); i += step {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunkText := r.cleanText(string(data[i:end]))
		if len(strings.TrimSpace(chunkText)) == 0 {
			continue
		}

		if err := encoder.Encode(Chunk{Text: chunkText}); err != nil {
			return fmt.Errorf("failed to write chunk: %w", err)
		}
	}

	return nil
}
