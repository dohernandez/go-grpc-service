#!/usr/bin/env bash

# Check if exclude-coverage.txt exists
if [ ! -f ./integration-exclude.coverpkg ]; then
  exit 0
fi

# Process each non-empty, non-comment line in exclude-coverage.txt
grep -v -E '^\s*#|^\s*$' ./integration-exclude.coverpkg | while read -r p || [ -n "$p" ]; do
  sed -i '' "/${p//\//\\/}/d" ./integration.coverprofile
done
