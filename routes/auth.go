package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"snippetier/auth"
	"snippetier/configs"
	"snippetier/db"
)

func setupAuthRoutes(g *echo.Group, _ *db.Storage, config *configs.Config) {
	g.GET("/login", loginHandler)
	g.GET("/login/github", githubLoginHandler(config.GithubClientId))
	g.GET("/github/callback", githubCallbackHandler(config))
}

func loginHandler(c echo.Context) error {
	return c.HTML(http.StatusOK, `<a href="/auth/login/github">Login with Github</a>`)
}

func githubLoginHandler(githubClientID string) echo.HandlerFunc {
	return func(c echo.Context) error {

		// Create the dynamic redirect URL for login
		redirectURL := fmt.Sprintf(
			"https://github.com/login/oauth/authorize?client_id=%s&scope=user:email",
			githubClientID,
		)

		return c.Redirect(301, redirectURL)
	}

}

func githubCallbackHandler(config *configs.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.QueryParam("code")
		if code == "" {
			return c.String(http.StatusBadRequest, "Missing code")
		}

		githubAccessToken := auth.GetGithubAccessToken(code, config)

		githubData := auth.GetGithubData(githubAccessToken)
		if githubData == "" {
			// Unauthorized users get an unauthorized message
			return c.String(http.StatusUnauthorized, "UNAUTHORIZED!")
		}

		// Prettifying the json
		var prettyJSON bytes.Buffer
		parser := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
		if parser != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "JSON parse error"})
		}

		// Return the prettified JSON as a string
		return c.String(http.StatusOK, prettyJSON.String())

	}
}
