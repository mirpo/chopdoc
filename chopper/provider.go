package chopper

import (
	"bufio"
	"fmt"

	"github.com/mirpo/chopdoc/config"
)

type Chunk struct {
	Text string `json:"chunk"`
}

type ChopperProvider interface {
	Chop() error
}

func NewChopper(chunkMethod config.ChunkMethod, cfg *config.Config, rw *bufio.ReadWriter) (ChopperProvider, error) {
	switch chunkMethod {
	case config.Char:
		return NewCharChopper(cfg, rw), nil
	case config.Word:
		return NewWordChopper(cfg, rw), nil
	case config.Sentence:
		return NewSentenceChopper(cfg, rw), nil
	default:
		return nil, fmt.Errorf("unsupported chunkMethod: %s", chunkMethod)
	}
}
