package chopper

import (
	"bufio"
	"encoding/json"
	"strings"

	"github.com/mirpo/chopdoc/config"
)

type CharChopper struct {
	BaseChopper
}

func NewCharChopper(cfg *config.Config, rw *bufio.ReadWriter) *CharChopper {
	scanner := bufio.NewScanner(rw.Reader)
	scanner.Split(bufio.ScanRunes)

	return &CharChopper{
		BaseChopper: BaseChopper{
			cfg:     cfg,
			encoder: json.NewEncoder(rw.Writer),
			scanner: scanner,
		},
	}
}

func (c *CharChopper) scanInput() error {
	var builder strings.Builder
	builder.Grow(c.cfg.ChunkSize)
	step := c.cfg.ChunkSize - c.cfg.Overlap

	for c.scanner.Scan() {
		builder.WriteString(c.scanner.Text())

		if builder.Len() >= c.cfg.ChunkSize {
			chunk := builder.String()
			if err := c.writeChunk(chunk); err != nil {
				return err
			}

			if step > len(chunk) {
				builder.Reset()
			} else {
				builder.Reset()
				builder.WriteString(chunk[step:])
			}
		}
	}

	if builder.Len() > 0 {
		return c.writeChunk(builder.String())
	}

	return c.scanner.Err()
}

func (c *CharChopper) Chop() error {
	return c.scanInput()
}
