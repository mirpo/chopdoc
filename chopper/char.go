package chopper

import (
	"bufio"
	"encoding/json"

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
	chunk := ""
	step := c.cfg.ChunkSize - c.cfg.Overlap

	for c.scanner.Scan() {
		chunk += c.scanner.Text()

		if len(chunk) >= c.cfg.ChunkSize {
			if err := c.writeChunk(chunk); err != nil {
				return err
			}

			if step > len(chunk) {
				chunk = ""
			} else {
				chunk = chunk[step:]
			}
		}
	}

	if len(chunk) > 0 {
		return c.writeChunk(chunk)
	}

	return c.scanner.Err()
}

func (c *CharChopper) Chop() error {
	return c.scanInput()
}
