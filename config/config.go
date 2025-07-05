package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
)

type ChunkMethod string

const (
	Char      ChunkMethod = "char"
	Word      ChunkMethod = "word"
	Sentence  ChunkMethod = "sentence"
	Recursive ChunkMethod = "recursive"
	Markdown  ChunkMethod = "markdown"
)

type CleaningMode string

const (
	CleanNormal     CleaningMode = "normal"
	CleanAggressive CleaningMode = "aggressive"
	CleanTrim       CleaningMode = "trim"
	CleanNone       CleaningMode = "none"
)

type Config struct {
	InputFile      string
	OutputFile     string
	Method         ChunkMethod
	ChunkSize      int
	Overlap        int
	CleaningMode   CleaningMode
	Piped          bool
	MarkdownHeader string
	MarkdownLevels []int
	StripHeaders   bool
	AddMetadata    bool
}

func NewConfig() *Config {
	return &Config{
		ChunkSize:      1000,
		Overlap:        0,
		CleaningMode:   CleanNone,
		Piped:          false,
		MarkdownHeader: "1-6",
		MarkdownLevels: []int{1, 2, 3, 4, 5, 6},
		StripHeaders:   false,
		AddMetadata:    false,
	}
}

func (c *Config) ParseMarkdownHeader() error {
	re := regexp.MustCompile(`^([1-6])-([1-6])$`)
	matches := re.FindStringSubmatch(c.MarkdownHeader)

	if matches == nil {
		return fmt.Errorf("invalid markdown header format: %s, expected format like '1-6'", c.MarkdownHeader)
	}

	start, _ := strconv.Atoi(matches[1])
	end, _ := strconv.Atoi(matches[2])

	if start > end {
		return fmt.Errorf("start level (%d) must be less than or equal to end level (%d)", start, end)
	}

	c.MarkdownLevels = make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		c.MarkdownLevels = append(c.MarkdownLevels, i)
	}

	return nil
}

func (c *Config) Validate() error {
	if !c.Piped {
		if c.InputFile == "" {
			return fmt.Errorf("input file is required")
		}
	}

	if c.OutputFile != "" && filepath.Ext(c.OutputFile) != ".jsonl" {
		return fmt.Errorf("output file must have .jsonl extension")
	}

	if c.ChunkSize <= 0 {
		return fmt.Errorf("chunk size must be greater than 0")
	}

	if c.Overlap >= c.ChunkSize {
		return fmt.Errorf("overlap must be less than chunk size")
	}

	validMethods := map[ChunkMethod]bool{
		Char:      true,
		Word:      true,
		Sentence:  true,
		Recursive: true,
		Markdown:  true,
	}
	if !validMethods[c.Method] {
		return fmt.Errorf("invalid chunking method: '%s'", c.Method)
	}

	if c.Method == Recursive && c.Overlap != 0 {
		fmt.Printf("warning: currently Recursive chopper doesn't support overlap, setting overlap to 0\n")
		c.Overlap = 0
	}

	if c.Method == Markdown {
		if err := c.ParseMarkdownHeader(); err != nil {
			return err
		}
	}

	return nil
}
