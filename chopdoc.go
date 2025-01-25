package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mirpo/chopdoc/config"
	"github.com/mirpo/chopdoc/runner"
)

type Stats struct {
	Type         string    `json:"type"`
	TotalChars   int       `json:"total_chars"`
	ChunkSize    int       `json:"chunk_size"`
	Overlap      int       `json:"overlap"`
	CleaningMode string    `json:"cleaning_mode"`
	StartedAt    time.Time `json:"started_at"`
	CompletedAt  time.Time `json:"completed_at"`
}

func main() {
	cfg := config.NewConfig()

	flag.StringVar(&cfg.InputFile, "input", "", "Input file path")
	flag.StringVar(&cfg.OutputFile, "output", "", "Output file path (must end with .jsonl)")
	flag.IntVar(&cfg.ChunkSize, "size", 1000, "Chunk size in characters")
	flag.IntVar(&cfg.Overlap, "overlap", 0, "Overlap size in characters")
	clean := flag.String("clean", "none", "Cleaning mode: none, normal, aggressive")
	flag.Parse()

	cfg.CleaningMode = config.CleaningMode(*clean)

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	r := runner.NewRunner(cfg)
	if err := r.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
