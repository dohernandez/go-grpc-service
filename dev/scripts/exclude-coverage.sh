#!/bin/sh

# Check if exclude-coverage.txt exists
if [ ! -f ./exclude-coverage.txt ]; then
  exit 0
fi

# Process each line in exclude-coverage.txt
while read -r p || [ -n "$p" ]; do
  sed -i '' "/${p//\//\\/}/d" ./coverage.out
done < ./exclude-coverage.txt
