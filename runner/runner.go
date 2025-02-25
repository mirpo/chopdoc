package runner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

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
		absPath, err := filepath.Abs(r.cfg.InputFile)
		if err != nil {
			return err
		}
		input, err = os.Open(absPath)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
	}
	defer input.Close()

	var output *os.File
	if r.cfg.OutputFile != "" {
		absPath, err := filepath.Abs(r.cfg.OutputFile)
		if err != nil {
			return err
		}
		output, err = os.Create(absPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
	} else {
		output = os.Stdout
	}
	defer output.Close()

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

	err = rw.Writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush buffers: %w", err)
	}

	return nil
}
