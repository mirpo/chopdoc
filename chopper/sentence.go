package chopper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/mirpo/chopdoc/config"
)

type SentenceChopper struct {
	BaseChopper
}

func NewSentenceChopper(cfg *config.Config, rw *bufio.ReadWriter) *SentenceChopper {
	scanner := bufio.NewScanner(rw.Reader)
	scanner.Split(scanSentences)

	return &SentenceChopper{
		BaseChopper: BaseChopper{
			cfg:     cfg,
			encoder: json.NewEncoder(rw.Writer),
			scanner: scanner,
		},
	}
}

func scanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	endPunctuation := regexp.MustCompile(`([.!?]+)(\s*)`)
	loc := endPunctuation.FindSubmatchIndex(data)

	if loc != nil {
		endIdx := loc[1]
		return endIdx, bytes.TrimRight(data[:endIdx], " \n\r\t"), nil
	}

	if atEOF {
		return len(data), bytes.TrimRight(data, " \n\r\t"), nil
	}

	return 0, nil, nil
}

func (s *SentenceChopper) scanInput() error {
	var sentences []string

	for s.scanner.Scan() {
		sentences = append(sentences, s.scanner.Text())

		if len(sentences) >= s.cfg.ChunkSize {
			chunk := strings.Join(sentences, " ")

			if err := s.writeChunk(chunk); err != nil {
				return err
			}

			if len(sentences) > s.cfg.Overlap {
				sentences = sentences[len(sentences)-s.cfg.Overlap:]
			} else {
				sentences = nil
			}
		}
	}

	if len(sentences) > 0 {
		chunk := strings.Join(sentences, " ")
		if err := s.writeChunk(chunk); err != nil {
			return err
		}
	}

	return s.scanner.Err()
}

func (s *SentenceChopper) Chop() error {
	return s.scanInput()
}
