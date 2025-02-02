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
	encoder := json.NewEncoder(rw.Writer)
	scanner := bufio.NewScanner(rw.Reader)
	scanner.Split(bufio.ScanBytes)

	return &CharChopper{
		BaseChopper: BaseChopper{
			cfg:     cfg,
			encoder: encoder,
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
			chunk = chunk[step:]
		}
	}

	if len(chunk) > 0 {
		return c.writeChunk(chunk)
	}

	return c.scanner.Err()
}

func (с *CharChopper) Chop() error {
	return с.scanInput()
}
