#!/bin/bash

FILE1=$1
FILE2=$2

if [ ! -f "$FILE1" ] || [ ! -f "$FILE2" ]; then
    echo "Error: One or both files do not exist."
    exit 1
fi

diff_output=$(diff "$FILE1" "$FILE2")

if [ $? -eq 0 ]; then
    echo "Files are identical."
else
    echo "Files differ:"
    echo "$diff_output"
    exit 1
fi
