#!/bin/bash

echo "🔴 Running BAD input tests (expecting ERROR):"
for file in bad_inputs/*.txt; do
  echo -n "Testing $file... "
  output=$(go run . "$file")
  if [ "$output" == "ERROR" ]; then
    echo "✅ Passed"
  else
    echo "❌ Failed (got: $output)"
  fi
done

echo
echo "🟢 Running GOOD input tests (expecting valid output):"
for file in good_inputs/*.txt; do
  echo "Testing $file:"
  go run . "$file"
  echo
done
