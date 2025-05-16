package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func DisambiguationHandler(title string, root string, paths []Route) echo.HandlerFunc {
	return func(c echo.Context) error {
		links := make([]string, len(paths))
		for i, path := range paths {
			links[i] = strings.TrimSuffix(path.Actual[len(root):], ".md")
		}

		data := map[string]any{
			"Title": title,
			"Links": links,
		}
		return c.Render(http.StatusOK, "disambiguation.html", data)
	}
}
