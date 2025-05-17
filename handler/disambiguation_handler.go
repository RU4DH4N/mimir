package handler

import (
	"net/http"
	"strings"

	"github.com/RU4DH4N/mimir/helper"
	"github.com/labstack/echo/v4"
)

func DisambiguationHandler(title string, root string, pages []PageData) echo.HandlerFunc {
	return func(c echo.Context) error {
		links := make([]string, len(pages))
		for i, page := range pages {
			links[i] = helper.Linkify(strings.TrimSuffix(page.Path[len(root)+1:], ".md"))
		}

		data := map[string]any{
			"Title": title,
			"Links": links,
		}
		return c.Render(http.StatusOK, "disambiguation.html", data)
	}
}
