#!/bin/sh

# Check if exclude-coverage.txt exists
if [ ! -f ./exclude-coverage.txt ]; then
  echo "exclude-coverage.txt not found. Exiting."
  exit 1
fi

# Process each line in exclude-coverage.txt
while read -r p || [ -n "$p" ]; do
  sed -i '' "/${p//\//\\/}/d" ./coverage.out
done < ./exclude-coverage.txt
