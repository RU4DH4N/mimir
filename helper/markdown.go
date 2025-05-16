package helper

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/toc"
)

func RenderMarkdown(path string, table bool) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve markdown path: %w", err)
	}

	content, err := os.ReadFile(abs)
	if err != nil {
		return "", fmt.Errorf("failed to read markdown file: %w", err)
	}

	// Prepare extensions based on 'table' flag
	extensions := []goldmark.Extender{extension.GFM}
	if table {
		extensions = append(extensions, &toc.Extender{})
	}

	md := goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown to HTML: %w", err)
	}

	return buf.String(), nil
}
