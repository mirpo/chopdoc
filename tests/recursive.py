import json
import argparse
from langchain.text_splitter import RecursiveCharacterTextSplitter

SEPARATORS = ["\n\n", "\n", " ", ".", ",", ""]

parser = argparse.ArgumentParser(description="Split text from input.txt and save to output.jsonl")
parser.add_argument("--size", type=int, default=60, help="Chunk size (default 60)")
parser.add_argument("--overlap", type=int, default=10, help="Overlap chunks (default 10)")
parser.add_argument("--input", type=str, default="input.txt", help="Input filename  (default input.txt)")
parser.add_argument("--output", type=str, default="output.jsonl", help="Output filename (default output.jsonl)")
args = parser.parse_args()

with open(args.input, "r", encoding="utf-8") as file:
    text = file.read()

splitter = RecursiveCharacterTextSplitter(
    chunk_size=args.size,
    chunk_overlap=args.overlap,
    separators=["\n\n", "\n", " ", ".", ",", ""]
)

chunks = splitter.split_text(text)

with open(args.output, "w", encoding="utf-8") as output_file:
    for chunk in chunks:
        json.dump({"chunk": chunk}, output_file, ensure_ascii=False, separators=(',', ':'))
        output_file.write("\n")

print(f"Splitted {len(chunks)} chunks")
