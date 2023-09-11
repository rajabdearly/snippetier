package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"snippetier/db"
	"strconv"
)

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, braaaat!")
}

func GetAllSnippets(storage *db.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		snippets, _ := storage.GetAllSnippets()
		return c.JSON(http.StatusOK, snippets)
	}
}

func SaveSnippet(storage *db.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		var snippet db.Snippet
		if err := c.Bind(&snippet); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		savedSnippet, err := storage.SaveSnippet(snippet.Name, snippet.Description, snippet.Content)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save snippet"})
		}

		return c.JSON(http.StatusCreated, savedSnippet)
	}
}

func UpdateSnippet(storage *db.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		snippetID := c.Param("id")
		id, err := strconv.Atoi(snippetID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid snippet ID"})
		}

		var snippet db.Snippet
		if err := c.Bind(&snippet); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		updatedSnippet, err := storage.UpdateSnippet(id, snippet.Name, snippet.Description, snippet.Content)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update snippet"})
		}

		return c.JSON(http.StatusOK, updatedSnippet)
	}
}

func DeleteSnippet(storage *db.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		snippetID := c.Param("id")
		id, err := strconv.Atoi(snippetID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid snippet ID"})
		}

		err = storage.DeleteSnippet(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete snippet"})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
