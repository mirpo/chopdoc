package chopper

import (
	"bufio"
	"encoding/json"
	"strings"

	"github.com/mirpo/chopdoc/config"
)

var defaultSeparators = []string{"\n\n", "\n", " ", ".", ",", ""}

type RecursiveChopper struct {
	BaseChopper
	buffer strings.Builder
}

func NewRecursiveChopper(cfg *config.Config, rw *bufio.ReadWriter) *RecursiveChopper {
	scanner := bufio.NewScanner(rw.Reader)
	scanner.Split(bufio.ScanLines)

	return &RecursiveChopper{
		BaseChopper: BaseChopper{
			cfg:     cfg,
			encoder: json.NewEncoder(rw.Writer),
			scanner: scanner,
		},
	}
}

func (r *RecursiveChopper) scanInput() error {
	for r.scanner.Scan() {
		line := r.scanner.Text()
		r.buffer.WriteString(line + "\n")

		if r.buffer.Len() >= r.cfg.ChunkSize {
			if err := r.processBuffer(); err != nil {
				return err
			}
		}
	}

	if r.buffer.Len() > 0 {
		return r.processBuffer()
	}

	return r.scanner.Err()
}

func (r *RecursiveChopper) processBuffer() error {
	text := r.buffer.String()
	r.buffer.Reset()

	for len(text) > 0 {
		chunk, remaining, ok := r.splitText(text)
		if !ok {
			chunk = text[:r.cfg.ChunkSize]
			remaining = text[r.cfg.ChunkSize:]
		}
		if err := r.writeChunk(chunk); err != nil {
			return err
		}
		text = remaining
	}
	return nil
}

func (r *RecursiveChopper) splitText(text string) (chunk string, remaining string, ok bool) {
	if len(text) <= r.cfg.ChunkSize {
		return text, "", true
	}

	for _, sep := range defaultSeparators {
		piece := string([]rune(text)[:r.cfg.ChunkSize])
		if pos := strings.LastIndex(piece, sep); pos != -1 {
			return text[:pos+len(sep)], text[pos+len(sep):], true
		}
	}

	return "", text, false
}

func (r *RecursiveChopper) Chop() error {
	return r.scanInput()
}
