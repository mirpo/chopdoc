import json
import argparse
from langchain.text_splitter import MarkdownHeaderTextSplitter

parser = argparse.ArgumentParser(description="Split text from input.txt and save to output.jsonl")
parser.add_argument("--input", type=str, default="input.txt", help="Input filename  (default input.txt)")
parser.add_argument("--output", type=str, default="output.jsonl", help="Output filename (default output.jsonl)")
parser.add_argument("--strip-headers", action="store_false", dest="no_strip_headers", help="Strip headers (default: True)")
args = parser.parse_args()

with open(args.input, "r", encoding="utf-8") as file:
    text = file.read()

headers_to_split_on = [
    ("#", "Header 1"),
    ("##", "Header 2"),
    ("###", "Header 3"),
    ("####", "Header 4"),
    ("#####", "Header 5"),
    ("######", "Header 6"),
]

print(args)

splitter = MarkdownHeaderTextSplitter(
    headers_to_split_on=headers_to_split_on, strip_headers=args.no_strip_headers
)
chunks = splitter.split_text(text)

with open(args.output, "w", encoding="utf-8") as output_file:
    for chunk in chunks:
        json.dump({"chunk": chunk.page_content, "metadata": chunk.metadata}, output_file, ensure_ascii=False, separators=(',', ':'))
        output_file.write("\n")

print(f"Splitted {len(chunks)} chunks")
