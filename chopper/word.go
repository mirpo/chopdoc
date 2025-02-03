package chopper

import (
	"bufio"
	"encoding/json"
	"strings"

	"github.com/mirpo/chopdoc/config"
)

type WordChopper struct {
	BaseChopper
}

func NewWordChopper(cfg *config.Config, rw *bufio.ReadWriter) *WordChopper {
	scanner := bufio.NewScanner(rw.Reader)
	scanner.Split(bufio.ScanWords)

	return &WordChopper{
		BaseChopper: BaseChopper{
			cfg:     cfg,
			encoder: json.NewEncoder(rw.Writer),
			scanner: scanner,
		},
	}
}

func (w *WordChopper) scanInput() error {
	var words []string

	for w.scanner.Scan() {
		words = append(words, w.scanner.Text())

		if len(words) >= w.cfg.ChunkSize {
			chunk := strings.Join(words, " ")

			if err := w.writeChunk(chunk); err != nil {
				return err
			}

			if len(words) > w.cfg.Overlap {
				words = words[len(words)-w.cfg.Overlap:]
			} else {
				words = nil
			}
		}
	}

	if len(words) > 0 {
		chunk := strings.Join(words, " ")
		if err := w.writeChunk(chunk); err != nil {
			return err
		}
	}

	return w.scanner.Err()
}

func (w *WordChopper) Chop() error {
	return w.scanInput()
}
