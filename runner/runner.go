package runner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mirpo/chopdoc/chopper"
	"github.com/mirpo/chopdoc/config"
)

type Runner struct {
	cfg *config.Config
}

func NewRunner(cfg *config.Config) *Runner {
	return &Runner{
		cfg: cfg,
	}
}

func (r *Runner) Run() error {
	var input *os.File
	var err error

	if r.cfg.Piped {
		input = os.Stdin
	} else {
		if err := validatePath(r.cfg.InputFile); err != nil {
			return fmt.Errorf("invalid input file path: %w", err)
		}
		absPath, err := filepath.Abs(r.cfg.InputFile)
		if err != nil {
			return err
		}
		input, err = os.Open(absPath)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer input.Close()
	}

	var output *os.File
	if r.cfg.OutputFile != "" {
		if err := validatePath(r.cfg.OutputFile); err != nil {
			return fmt.Errorf("invalid output file path: %w", err)
		}
		absPath, err := filepath.Abs(r.cfg.OutputFile)
		if err != nil {
			return err
		}
		output, err = os.Create(absPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)
	rw := bufio.NewReadWriter(reader, writer)

	chopper, err := chopper.NewChopper(r.cfg.Method, r.cfg, rw)
	if err != nil {
		return fmt.Errorf("failed to create chopper: %w", err)
	}

	err = chopper.Chop()
	if err != nil {
		return fmt.Errorf("failed to chop file: %w", err)
	}

	err = rw.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush buffers: %w", err)
	}

	return nil
}

func validatePath(path string) error {
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal detected: %s", path)
	}
	return nil
}
