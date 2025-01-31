package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mirpo/chopdoc/config"
	"github.com/mirpo/chopdoc/runner"
)

func main() {
	cfg := config.NewConfig()

	flag.StringVar(&cfg.InputFile, "input", "", "Input file path")
	flag.StringVar(&cfg.OutputFile, "output", "", "Output file path (must end with .jsonl)")
	flag.IntVar(&cfg.ChunkSize, "size", 1000, "Chunk size in characters")
	flag.IntVar(&cfg.Overlap, "overlap", 0, "Overlap size in characters")
	method := flag.String("method", string(config.ByCharacters), "Chunking method: characters")
	clean := flag.String("clean", "none", "Cleaning mode: none, normal, aggressive")

	flag.Parse()

	cfg.CleaningMode = config.CleaningMode(*clean)
	cfg.Method = config.ChunkMethod(*method)

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
