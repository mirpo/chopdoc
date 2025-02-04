package config

import (
	"fmt"
	"path/filepath"
)

type ChunkMethod string

const (
	Char     ChunkMethod = "char"
	Word     ChunkMethod = "word"
	Sentence ChunkMethod = "sentence"
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
	Method       ChunkMethod
	ChunkSize    int
	Overlap      int
	CleaningMode CleaningMode
}

func NewConfig() *Config {
	return &Config{
		ChunkSize:    1000,
		Overlap:      0,
		CleaningMode: CleanNormal,
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
	validMethods := map[ChunkMethod]bool{
		Char:     true,
		Word:     true,
		Sentence: true,
	}
	if !validMethods[c.Method] {
		return fmt.Errorf("invalid method: %s", c.Method)
	}
	return nil
}
