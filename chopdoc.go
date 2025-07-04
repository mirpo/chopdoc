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
	flag.BoolVar(&ver, "version", false, "Get current version of chopdoc")
	flag.StringVar(&cfg.InputFile, "input", "", "Input file path")
	flag.StringVar(&cfg.OutputFile, "output", "", "Output file path (must end with .jsonl)")
	flag.IntVar(&cfg.ChunkSize, "size", 1000, "Chunk size in characters")
	flag.IntVar(&cfg.Overlap, "overlap", 0, "Overlap size in characters")
	method := flag.String("method", string(config.Char), "Default chunking method: char")
	clean := flag.String("clean", "none", "Cleaning mode: none, normal, aggressive")

	// used only in markdown chopper
	flag.StringVar(&cfg.MarkdownHeader, "headers", "1-6", "Header levels to use for markdown method (e.g. 1-6, 2-4)")
	flag.BoolVar(&cfg.StripHeaders, "strip-headers", false, "Remove headers from content (default false, markdown method only)")
	flag.BoolVar(&cfg.AddMetadata, "add-metadata", false, "Include header metadata in output (default false, markdown method only)")

	flag.Parse()

	if ver {
		slog.Info("chopdoc", "version", version, "commit", commit)
		return
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		slog.Error("failed to check stdin", "err", err)
		os.Exit(1)
	}
	cfg.Piped = (stat.Mode()&os.ModeCharDevice) == 0 && cfg.InputFile == ""
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
