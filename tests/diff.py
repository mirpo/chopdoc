import sys
import difflib

def compare_files(file1, file2):
    try:
        with open(file1, 'r', encoding='utf-8') as f1, open(file2, 'r', encoding='utf-8') as f2:
            lines1 = f1.readlines()
            lines2 = f2.readlines()
        
        diff = list(difflib.unified_diff(lines1, lines2, fromfile=file1, tofile=file2, lineterm=''))
        if not diff:
            sys.exit(0)
        else:
            print("Files differ:")
            print('\n'.join(diff))
            sys.exit(1)
    except FileNotFoundError as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python diff.py <file1> <file2>")
        sys.exit(1)
    
    compare_files(sys.argv[1], sys.argv[2])
