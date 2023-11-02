package main

import (
	"html/template"
	"log"
	"net/http"
	"snippetier/configs"
	"snippetier/db"
	"snippetier/routes"
	renderer "snippetier/templates"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config, err := configs.GetConfig()

	if err != nil {
		log.Fatal("Error while reading config: ", err)
	}

	storage, err := db.GetConnection()
	defer storage.CloseConnection()
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}

	t := &renderer.Template{
		Templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e := echo.New()
	e.Renderer = t
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	routes.SetupRoutes(e, storage, config)
	e.GET("/", rootHandler)

	err = e.Start(":1323")
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, braaaat!")
}
