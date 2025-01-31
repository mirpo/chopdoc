# chopdoc

A command-line tool for splitting documents into chunks, optimized for RAG (Retrieval-Augmented Generation) and LLM applications.

## Features
- Supports chunking text files into configurable sizes
- Configurable chunk size and overlap
- Text cleaning and normalization
- JSONL output format
- Support for plain text files (with Markdown, PDF planned)

## Installation

### Global Installation
```bash
git clone https://github.com/mirpo/chopdoc.git
cd chopdoc
make install
```

Or using Go:
```bash
go install github.com/mirpo/chopdoc@latest
```

### Local Build
```bash
git clone https://github.com/mirpo/chopdoc.git
cd chopdoc
go build
```

## Usage

```bash
chopdoc -input pg_essay.txt -output chunks.jsonl -size 1000
chopdoc -input pg_essay.txt -output chunks.jsonl -size 1000 -overlap 100
chopdoc -input pg_essay.txt -output chunks.jsonl -size 1000 -overlap 100 -clean aggressive -strategy char
```

### Options

| Option      | Description                                 | Default  |
| ----------- | ------------------------------------------- | -------- |
| `-input`    | Input file path                             | Required |
| `-output`   | Output file path (.jsonl)                   | Required |
| `-size`     | Chunk size (characters/words/sentences)     | 1000     |
| `-overlap`  | Overlap between chunks                      | 0        |
| `-clean`    | Cleaning mode (none/normal/aggressive)      | none     |
| `-strategy` | Chunking strategy (character/word/sentence) | char     |

### Output Format

Each chunk is written as a JSON line:
```json
{"chunk": "content here"}
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Run tests: `go test ./...`
4. Submit a pull request

## License

MIT
