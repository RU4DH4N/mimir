package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/RU4DH4N/mimir/helper"
	"github.com/labstack/echo/v4"
)

type Link struct {
	URL  string
	Text string
}

func DisambiguationHandler(title, root, route string, pages []PageData) echo.HandlerFunc {
	return func(c echo.Context) error {
		links := make([]Link, len(pages))
		for i, page := range pages {
			links[i] = Link{
				Text: helper.Linkify(strings.TrimSuffix(page.Path[len(root)+1:], ".md")),
				URL:  route + "/v" + strconv.Itoa(i+1),
			}
		}

		data := map[string]any{
			"Title": title,
			"Links": links,
		}
		return c.Render(http.StatusOK, "disambiguation.html", data)
	}
}
