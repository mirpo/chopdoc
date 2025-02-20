import json
from langchain.text_splitter import RecursiveCharacterTextSplitter

CHUNK_SIZE = 60
CHUNK_OVERLAP = 10
SEPARATORS = ["\n\n", "\n", " ", ".", ",", ""]

with open("input.txt", "r", encoding="utf-8") as file:
    text = file.read()

splitter = RecursiveCharacterTextSplitter(
    chunk_size=CHUNK_SIZE,
    chunk_overlap=CHUNK_OVERLAP,
    separators=SEPARATORS
)

chunks = splitter.split_text(text)

with open("output.jsonl", "w", encoding="utf-8") as output_file:
    for chunk in chunks:
        json.dump({"chunk": chunk}, output_file, ensure_ascii=False)
        output_file.write("\n")

print(f"✅ Splitted {len(chunks)} chunks")
