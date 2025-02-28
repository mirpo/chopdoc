package chopper

import (
	"bufio"
	"encoding/json"

	"github.com/mirpo/chopdoc/config"
)

type MarkdownChopper struct {
	BaseChopper
}

func NewMarkdownChopper(cfg *config.Config, rw *bufio.ReadWriter) *MarkdownChopper {
	scanner := bufio.NewScanner(rw.Reader)
	scanner.Split(bufio.ScanRunes)

	return &MarkdownChopper{
		BaseChopper: BaseChopper{
			cfg:     cfg,
			encoder: json.NewEncoder(rw.Writer),
			scanner: scanner,
		},
	}
}

func (m *MarkdownChopper) scanInput() error {
	chunk := ""
	step := m.cfg.ChunkSize - m.cfg.Overlap

	for m.scanner.Scan() {
		chunk += m.scanner.Text()

		if len(chunk) >= m.cfg.ChunkSize {
			if err := m.writeChunk(chunk); err != nil {
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
		return m.writeChunk(chunk)
	}

	return m.scanner.Err()
}

func (m *MarkdownChopper) Chop() error {
	return m.scanInput()
}
