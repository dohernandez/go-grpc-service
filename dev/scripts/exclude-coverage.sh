#!/usr/bin/env bash

# Check if exclude-coverage.txt exists
if [ ! -f ./integration-exclude.coverpkg ]; then
  exit 0
fi

# Check if the system is macOS or Linux
OS_TYPE=$(uname -s)

# Process each non-empty, non-comment line in exclude-coverage.txt
grep -v -E '^\s*#|^\s*$' ./integration-exclude.coverpkg | while read -r p || [ -n "$p" ]; do
  # Escape special characters for sed
  ESCAPED_PATH=$(echo "$p" | sed 's/[&/\]/\\&/g')

  # On macOS, use -i '' for in-place editing; on Linux, use -i without any argument
  if [[ "$OS_TYPE" == "Darwin" ]]; then
    # macOS: use -i '' for in-place editing
    sed -i '' "/$ESCAPED_PATH/d" ./integration.coverprofile
  else
    # Linux: use -i without any argument for in-place editing
    sed -i "/$ESCAPED_PATH/d" ./integration.coverprofile
  fi
done
