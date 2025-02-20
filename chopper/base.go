package chopper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mirpo/chopdoc/cleaner"
	"github.com/mirpo/chopdoc/config"
)

type BaseChopper struct {
	cfg     *config.Config
	encoder *json.Encoder
	scanner *bufio.Scanner
}

func (b *BaseChopper) cleanChunk(chunk string) string {
	return cleaner.Clean(chunk, b.cfg.CleaningMode)
}

func (b *BaseChopper) writeChunk(chunk string) error {
	b.encoder.SetEscapeHTML(false)
	chunk = b.cleanChunk(chunk)

	if len(strings.TrimSpace(chunk)) == 0 {
		return nil
	}

	if err := b.encoder.Encode(Chunk{Text: chunk}); err != nil {
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	return nil
}
