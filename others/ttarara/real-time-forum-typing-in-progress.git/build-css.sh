#!/bin/bash
# Build Tailwind CSS for production

echo "Building Tailwind CSS..."
npx tailwindcss -i ./frontend/src/input.css -o ./frontend/dist/output.css --minify

echo "CSS build complete!"

