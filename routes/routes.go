package routes

import (
	"net/http"
	"snippetier/db"

	"github.com/labstack/echo/v4"
)

const userIdHeader = "sn-trusted-user-id"

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, braaaat!")
}

func SetupRoutes(e *echo.Echo, s *db.Storage) {

	apiGroup := e.Group("api")

	apiGroup.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Test")
	})

	snippetsGroup := apiGroup.Group("/snippets")
	SetupSnippetsRoutes(snippetsGroup, s)

	usersGroup := apiGroup.Group("/users")
	SetupUserRoutes(usersGroup, s)
}
