#!/usr/bin/env sh

set -e

INPUT="./wiki-example/tailwind/input.css"

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "❌ Required command '$1' is not installed or not in PATH."
    exit 1
  fi
}

require_command npm
require_command npx
require_command go

grep -E '^@(import|plugin) ' "$INPUT" | while read -r line; do
  pkg=$(echo "$line" | sed -E 's/^@(import|plugin) +["'\'']([^"'\'';]+)["'\''];?/\2/')

  if [[ "$pkg" == .* || "$pkg" == /* ]]; then
    continue
  fi

  if ! npm ls "$pkg" --depth=0 >/dev/null 2>&1; then
    echo "❌ Package '$pkg' is not installed."
    echo "  Run: npm install $pkg"
    exit 1
  fi
done

npx @tailwindcss/cli -i "$INPUT" -o "./wiki-example/static/css/tailwind.css"

go run .
