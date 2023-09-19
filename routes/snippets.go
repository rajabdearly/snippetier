package routes

import (
	"net/http"
	"snippetier/db"
	"snippetier/db/repo"
	"strconv"

	"github.com/labstack/echo/v4"
)

func SetupSnippetsRoutes(g *echo.Group, storage *db.Storage) {
	g.GET("", getAllSnippets(storage))
	g.POST("/new", saveSnippet(storage))
	g.PUT("/:id", updateSnippet(storage))
	g.DELETE("/:id", deleteSnippet(storage))
}

func getAllSnippets(storage *db.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		snippets, _ := storage.SnippetsRepo.GetAllSnippets()
		return c.JSON(http.StatusOK, snippets)
	}
}

func saveSnippet(storage *db.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		parsedUserId := c.Request().Header.Get(userIdHeader)
		userId, err := strconv.Atoi(parsedUserId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		var snippet repo.Snippet
		if err := c.Bind(&snippet); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		savedSnippet, err := storage.SnippetsRepo.SaveSnippet(userId, snippet.Name, snippet.Description, snippet.Content)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save snippet"})
		}

		return c.JSON(http.StatusCreated, savedSnippet)
	}
}

func updateSnippet(storage *db.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		parsedUserId := c.Request().Header.Get(userIdHeader)
		userId, err := strconv.Atoi(parsedUserId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}
		parsedSnippetID := c.Param("id")
		snippetID, err := strconv.Atoi(parsedSnippetID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid snippet ID"})
		}

		var snippet repo.Snippet
		if err := c.Bind(&snippet); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		updatedSnippet, err := storage.SnippetsRepo.UpdateSnippet(userId, snippetID, snippet.Name, snippet.Description, snippet.Content)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update snippet"})
		}

		return c.JSON(http.StatusOK, updatedSnippet)
	}
}

func deleteSnippet(storage *db.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		snippetID := c.Param("id")
		id, err := strconv.Atoi(snippetID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid snippet ID"})
		}

		err = storage.SnippetsRepo.DeleteSnippet(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete snippet"})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
