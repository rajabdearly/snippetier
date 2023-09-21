package routes

import (
	"net/http"
	"snippetier/db"
	"snippetier/db/repo"
	"strconv"

	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(g *echo.Group, s *db.Storage) {
	g.GET("/me", getUserMe(s))
	g.GET("/:id", getUserById(s))
	g.PUT("/:id", updateUser(s))
}

// getUserById retrieves a user by ID and returns it.
func getUserById(storage *db.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("id")
		id, err := strconv.Atoi(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		user, err := storage.UsersRepo.GetUserByID(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user"})
		}

		return c.JSON(http.StatusOK, user)
	}
}

func getUserMe(storage *db.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().Header.Get(UserIdHeader)
		id, err := strconv.Atoi(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		user, err := storage.UsersRepo.GetUserByID(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user"})
		}

		return c.JSON(http.StatusOK, user)
	}
}

// updateUser updates an existing user and returns the updated user.
func updateUser(storage *db.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("id")
		id, err := strconv.Atoi(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		var user repo.User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		updatedUser, err := storage.UsersRepo.UpdateUser(id, user.Username, user.Email, user.FullName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}
