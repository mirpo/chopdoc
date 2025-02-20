# chopdoc

A command-line tool for splitting documents into chunks, optimized for RAG (Retrieval-Augmented Generation) and LLM applications.

## Features
- Supports chunking methods: characters, words, sentences, recursive
- Configurable chunk size and overlap
- Text cleaning and normalization
- JSONL output format
- Supported formats: txt (plain test)

## Installation

[Homebrew](https://brew.sh/):
```shell
brew tap mirpo/homebrew-tools
brew install chopdoc
```

Using `go install`:
```shell
go install github.com/mirpo/chopdoc@latest
```

### Local Build
```shell
git clone https://github.com/mirpo/chopdoc.git
cd chopdoc
make build
```

## Usage

```bash
chopdoc -input pg_essay.txt -output chunks.jsonl -size 1000 -clean aggressive
chopdoc -input pg_essay.txt -output chunks.jsonl -size 1000 -overlap 100
chopdoc -input pg_essay.txt -output chunks.jsonl -size 1000 -overlap 100 -method char -clean aggressive
chopdoc -input pg_essay.txt -output chunks.jsonl -size 1000 -overlap 100 -method word
chopdoc -input pg_essay.txt -output chunks.jsonl -size 10   -overlap 1   -method sentence
chopdoc -input pg_essay.txt -output chunks.jsonl -size 100  -overlap 0   -method recursive
```

chopdoc can be piped:
```bash
cat pg_essay.txt | chopdoc -size 1 -method sentence
cat pg_essay.txt | chopdoc -size 1 -method sentence > piped.jsonl
cat pg_essay.txt | chopdoc -size 1 -method sentence -output output_as_arg.jsonl
```

### Options

| Option     | Description                            | Default  |
| ---------- | -------------------------------------- | -------- |
| `-input`   | Input file path                        | Required |
| `-output`  | Output file path (.jsonl)              | Required |
| `-size`    | Chunk size                             | 1000     |
| `-overlap` | Overlap between chunks                 | 0        |
| `-clean`   | Cleaning mode (none/normal/aggressive) | none     |
| `-method`  | Chunking method                        | char     |

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
