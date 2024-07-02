#!/bin/bash

test_dirs=$(find . -type f -name '*_test.go' -exec dirname {} \; | sort -u)
coverage_file="coverage.out"
echo "mode: atomic" > "$coverage_file"
for dir in $test_dirs; do
    temp_coverage_file="$dir/coverage.out"
    go test "$dir" -coverprofile="$temp_coverage_file" -covermode=atomic -v
    if [ -f "$temp_coverage_file" ]; then
        tail -n +2 "$temp_coverage_file" >> "$coverage_file"
        rm "$temp_coverage_file"
    fi
done

go tool cover -html="$coverage_file" -o coverage.html
open coverage.html
