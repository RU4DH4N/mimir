package handler

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// this is terrible I know, I'll fix it in a bit
func DisambiguationHandler(title string, paths []string) echo.HandlerFunc {
	return func(c echo.Context) error {
		pathParts := make([][]string, len(paths))
		for i, p := range paths {
			abs := filepath.ToSlash(p)
			parts := strings.Split(abs, "/")
			pathParts[i] = parts
		}

		links := make([]string, len(paths))
		for i := range paths {
			suffixLen := 1
			for {
				conflict := false
				thisSuffix := joinLast(pathParts[i], suffixLen)

				for j, _ := range paths {
					if i == j {

						continue
					}
					if joinLast(pathParts[j], suffixLen) == thisSuffix {
						conflict = true
						break
					}
				}
				if !conflict {
					links[i] = thisSuffix
					break
				}
				suffixLen++
				if suffixLen > len(pathParts[i]) {
					links[i] = strings.Join(pathParts[i], "/")
					break
				}
			}
		}

		data := map[string]any{
			"Title": title,
			"Links": links,
		}
		return c.Render(http.StatusOK, "disambiguation.html", data)
	}
}

func joinLast(parts []string, n int) string {
	if n > len(parts) {
		n = len(parts)
	}
	return strings.Join(parts[len(parts)-n:], "/")
}
