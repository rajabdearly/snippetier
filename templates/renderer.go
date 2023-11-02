package templates

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(
	w io.Writer,
	name string,
	data interface{},
	e echo.Context,
) error {
	e.Logger().Debug(e.Cookies())
	return t.Templates.ExecuteTemplate(w, name, data)
}
