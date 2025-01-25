package config

import (
	"fmt"
	"path/filepath"
)

type CleaningMode string

const (
	CleanNormal     CleaningMode = "normal"
	CleanAggressive CleaningMode = "aggressive"
	CleanNone       CleaningMode = "none"
)

type Config struct {
	InputFile    string
	OutputFile   string
	ChunkSize    int
	Overlap      int
	CleaningMode CleaningMode
	Stats        bool
}

func NewConfig() *Config {
	return &Config{
		ChunkSize:    1000,
		Overlap:      0,
		CleaningMode: CleanNormal,
		Stats:        false,
	}
}

func (c *Config) Validate() error {
	if c.InputFile == "" {
		return fmt.Errorf("input file is required")
	}
	if c.OutputFile == "" {
		return fmt.Errorf("output file is required")
	}
	if c.ChunkSize <= 0 {
		return fmt.Errorf("chunk size must be greater than 0")
	}
	if c.Overlap >= c.ChunkSize {
		return fmt.Errorf("overlap must be less than chunk size")
	}
	if filepath.Ext(c.OutputFile) != ".jsonl" {
		return fmt.Errorf("output file must have .jsonl extension")
	}
	return nil
}
