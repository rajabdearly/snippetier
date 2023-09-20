package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"snippetier/configs"
	"snippetier/db"
	"snippetier/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config, err := configs.GetConfig()

	if err != nil {
		log.Fatal("Error while reading config: ", err)
	}

	db.SetupNewTestDb(config.DbName)
	storage, err := db.GetConnection(config.DbName)
	defer storage.CloseConnection()
	if err != nil {
		log.Fatal("Seding failed with err: ", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	routes.SetupRoutes(e, storage)
	fmt.Println(config)

	e.GET("/", rootHandler)

	e.GET("/login/github/", githubLoginHandler(config.GithubClientId))

	// Github callback
	e.GET("/auth/github/callback", githubCallbackHandler(config))

	// Route where the authenticated user is redirected to
	e.GET("/loggedin", loggedinHandler)

	err = e.Start(":1323")
	if err != nil {
		log.Fatal(err)
	}
}

func loggedinHandler(c echo.Context) error {

	return c.HTML(http.StatusOK, `<h1> You are logged in! </h1>`)
	githubData := c.Request().Header.Get("githubData")

	if githubData == "" {
		// Unauthorized users get an unauthorized message
		return c.String(http.StatusUnauthorized, "UNAUTHORIZED!")
	}

	// Prettifying the json
	var prettyJSON bytes.Buffer
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "JSON parse error"})
	}

	// Return the prettified JSON as a string
	return c.String(http.StatusOK, prettyJSON.String())
}

func rootHandler(c echo.Context) error {
	return c.HTML(http.StatusOK, `<a href="/login/github/">LOGIN</a>`)
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

		githubAccessToken := getGithubAccessToken(code, config)

		githubData := getGithubData(githubAccessToken)
		if githubData == "" {
			// Unauthorized users get an unauthorized message
			return c.String(http.StatusUnauthorized, "UNAUTHORIZED!")
		}

		// Prettifying the json
		var prettyJSON bytes.Buffer
		parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
		if parserr != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "JSON parse error"})
		}

		// Return the prettified JSON as a string
		return c.String(http.StatusOK, prettyJSON.String())

	}
}

func getGithubAccessToken(code string, config *configs.Config) string {
	fmt.Println("Config: ", config)

	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     config.GithubClientId,
		"client_secret": config.GithubClientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Response body converted to stringified JSON
	respbody, _ := io.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken
}

func getGithubData(accessToken string) string {
	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Read the response as a byte slice
	respbody, _ := io.ReadAll(resp.Body)

	// Convert byte slice to string and return
	return string(respbody)
}
