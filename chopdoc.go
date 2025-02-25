package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/mirpo/chopdoc/config"
	"github.com/mirpo/chopdoc/runner"
)

var (
	version string = "dev-build"
	commit  string = "commit"
)

func main() {
	cfg := config.NewConfig()

	var ver bool
	flag.BoolVar(&ver, "version", false, "Get current version of sentences")
	flag.StringVar(&cfg.InputFile, "input", "", "Input file path")
	flag.StringVar(&cfg.OutputFile, "output", "", "Output file path (must end with .jsonl)")
	flag.IntVar(&cfg.ChunkSize, "size", 1000, "Chunk size in characters")
	flag.IntVar(&cfg.Overlap, "overlap", 0, "Overlap size in characters")
	method := flag.String("method", string(config.Char), "Default chunking method: char")
	clean := flag.String("clean", "none", "Cleaning mode: none, normal, aggressive")

	flag.Parse()

	if ver {
		slog.Info("chopdoc", "version", version, "commit", commit)
		return
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		slog.Error("failed to Stdin.Stat", "err", err)
		os.Exit(1)
	}
	cfg.Piped = (fi.Mode() & os.ModeNamedPipe) != 0
	cfg.CleaningMode = config.CleaningMode(*clean)
	cfg.Method = config.ChunkMethod(*method)

	if err := cfg.Validate(); err != nil {
		slog.Error("failed to validate config", "err", err)
		os.Exit(1)
	}

	r := runner.NewRunner(cfg)
	if err := r.Run(); err != nil {
		slog.Error("execution error", "err", err)
		os.Exit(1)
	}
}
