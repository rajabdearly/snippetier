package routes

import (
	"net/http"
	"snippetier/configs"
	"snippetier/db"

	"github.com/labstack/echo/v4"
)

const UserIdHeader = "sn-trusted-user-id"

// SetupRoutes sets up all the routes for the application
func SetupRoutes(e *echo.Echo, s *db.Storage, config *configs.Config) {

	apiGroup := e.Group("api")

	apiGroup.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Test")
	})

	snippetsGroup := apiGroup.Group("/snippets")
	SetupSnippetsRoutes(snippetsGroup, s)

	usersGroup := apiGroup.Group("/users")
	SetupUserRoutes(usersGroup, s)

	authGroup := e.Group("/auth")
	setupAuthRoutes(authGroup, s, config)
}
