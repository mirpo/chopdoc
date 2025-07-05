# chopdoc

A command-line tool for splitting documents into chunks, optimized for RAG (Retrieval-Augmented Generation) and LLM applications.

## Features
- Supports chunking methods: characters, words, sentences, recursive, markdown.
- Configurable chunk size and overlap
- Text cleaning and normalization
- JSONL output format
- Supported formats: txt (or any plain text)

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
chopdoc -input pg_essay.txt -output chunks.jsonl -size 100  -overlap 0   -method recursive
chopdoc -input pg_essay.txt -output chunks.jsonl                         -method markdown -strip-headers
chopdoc -input pg_essay.txt -output chunks.jsonl                         -method markdown -headers 1-2 -add-metadata
```

chopdoc can be piped:
```bash
cat pg_essay.txt | chopdoc -size 1 -method sentence
cat pg_essay.txt | chopdoc -size 1 -method sentence > piped.jsonl
cat pg_essay.txt | chopdoc -size 1 -method sentence -output output_as_arg.jsonl
```

### Options

```shell
  -add-metadata
        Include header metadata in output (default false, markdown method only)
  -clean string
        Cleaning mode: none, normal, aggressive (default "none")
  -headers string
        Header levels to use for markdown method (e.g. 1-6, 2-4) (default "1-6")
  -input string
        Input file path
  -method string
        Default chunking method: char (default "char")
  -output string
        Output file path (must end with .jsonl)
  -overlap int
        Overlap size in characters
  -size int
        Chunk size in characters (default 1000)
  -strip-headers
        Remove headers from content (default false, markdown method only)
  -version
        Get current version of chopdoc
```

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
