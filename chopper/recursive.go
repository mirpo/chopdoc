package chopper

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"

	"github.com/mirpo/chopdoc/config"
)

const (
	defaultBufferSize    = 64 * 1024
	bufferSizeMultiplier = 1
)

var (
	defaultSeparators = []string{"\n\n", "\n", " ", ".", ",", ""}
)

type RecursiveChopper struct {
	BaseChopper
	bufReader *bufio.Reader
	buffer    strings.Builder
}

func NewRecursiveChopper(cfg *config.Config, rw *bufio.ReadWriter) *RecursiveChopper {
	bufReader := bufio.NewReaderSize(rw.Reader, defaultBufferSize)

	return &RecursiveChopper{
		BaseChopper: BaseChopper{
			cfg:     cfg,
			encoder: json.NewEncoder(rw.Writer),
		},
		bufReader: bufReader,
	}
}

func (r *RecursiveChopper) scanInput() error {
	for {
		chunk, err := r.bufReader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		r.buffer.Write(chunk)

		if r.buffer.Len() >= r.cfg.ChunkSize*bufferSizeMultiplier || err == io.EOF {
			if err := r.processBuffer(); err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}
	}
	return nil
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
