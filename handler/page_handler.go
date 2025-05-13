package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/RU4DH4N/mimir/helper"
	"github.com/labstack/echo/v4"
)

func PageHandler(path string) echo.HandlerFunc {
	return func(c echo.Context) error {
		content, err := helper.RenderMarkdown(path)
		if err != nil {
			fmt.Printf("Leave the gun, take the '%s'.", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render page")
		}
		data := map[string]any{
			"Content": template.HTML(content),
		}
		return c.Render(http.StatusOK, "index.html", data)
	}
}
