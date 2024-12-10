#!/usr/bin/env bash

# Check if exclude-coverage.txt exists
if [ ! -f ./exclude-coverage.txt ]; then
  exit 0
fi

# Process each non-empty, non-comment line in exclude-coverage.txt
grep -v -E '^\s*#|^\s*$' ./exclude-coverage | while read -r p || [ -n "$p" ]; do
  sed -i '' "/${p//\//\\/}/d" ./integration.coverprofile
done
