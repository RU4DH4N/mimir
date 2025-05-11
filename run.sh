#!/bin/zsh

set -e

npx @tailwindcss/cli -i "./wiki-example/tailwind/input.css" -o "./wiki-example/static/css/tailwind.css"

go run main.go

