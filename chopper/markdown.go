package chopper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mirpo/chopdoc/config"
)

type header struct {
	Pattern string
	Name    string
	Level   int
}

type MarkdownChopper struct {
	BaseChopper
	headers   []header
	headerRgx *regexp.Regexp
	metadata  map[string]string
}

func createHeaders(levels []int) []header {
	headers := make([]header, len(levels))

	for i, level := range levels {
		headers[i] = header{
			Pattern: strings.Repeat("#", level),
			Name:    "Header " + strconv.Itoa(level),
			Level:   level,
		}
	}

	return headers
}

func createHeaderRegex(headers []header) *regexp.Regexp {
	patterns := make([]string, len(headers))
	for i, h := range headers {
		patterns[i] = "^(" + regexp.QuoteMeta(h.Pattern) + " )"
	}
	return regexp.MustCompile(strings.Join(patterns, "|"))
}

func NewMarkdownChopper(cfg *config.Config, rw *bufio.ReadWriter) *MarkdownChopper {
	headers := createHeaders(cfg.MarkdownLevels)

	return &MarkdownChopper{
		BaseChopper: BaseChopper{
			cfg:     cfg,
			encoder: json.NewEncoder(rw.Writer),
			scanner: bufio.NewScanner(rw.Reader),
		},
		headers:   headers,
		headerRgx: createHeaderRegex(headers),
		metadata:  make(map[string]string),
	}
}

func (m *MarkdownChopper) scanInput() error {
	var buffer strings.Builder

	for m.scanner.Scan() {
		line := m.scanner.Text()
		if len(line) == 0 {
			continue
		}

		if m.headerRgx.MatchString(line) {
			if buffer.Len() > 0 {
				if err := m.processBuffer(buffer.String()); err != nil {
					return err
				}
				buffer.Reset()
			}
			m.updateMetadata(line)

			if m.cfg.StripHeaders {
				continue
			}
		}

		buffer.WriteString(line + "\n")
	}

	if buffer.Len() > 0 {
		if err := m.processBuffer(buffer.String()); err != nil {
			return err
		}
	}

	return m.scanner.Err()
}

func (m *MarkdownChopper) updateMetadata(line string) {
	for _, header := range m.headers {
		headerPrefix := header.Pattern + " "
		if strings.HasPrefix(line, headerPrefix) {
			m.metadata[header.Name] = strings.TrimPrefix(line, headerPrefix)
			break
		}
	}
}

func (m *MarkdownChopper) processBuffer(chunk string) error {
	m.encoder.SetEscapeHTML(false)
	chunk = m.cleanChunk(chunk)

	if len(strings.TrimSpace(chunk)) == 0 {
		return nil
	}

	jsonlChunk := Chunk{Text: chunk}
	if m.cfg.AddMetadata {
		jsonlChunk.Metadata = m.metadata
	}

	if err := m.encoder.Encode(jsonlChunk); err != nil {
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	return nil
}

func (m *MarkdownChopper) Chop() error {
	return m.scanInput()
}
