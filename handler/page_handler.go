package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/RU4DH4N/mimir/helper"
	"github.com/labstack/echo/v4"
)

type PageData struct {
	Title string
	Path  string
	Table bool
}

func PageHandler(data PageData) echo.HandlerFunc {
	return func(c echo.Context) error {
		content, err := helper.RenderMarkdown(data.Path, data.Table)
		if err != nil {
			fmt.Printf("Leave the gun, take the '%s'.", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render page")
		}
		data := map[string]any{
			"Title":   data.Title,
			"Content": template.HTML(content),
		}
		return c.Render(http.StatusOK, "index.html", data)
	}
}
