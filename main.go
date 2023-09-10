package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"snippetier/db"
	"snippetier/routes"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", routes.Hello)

	db.New("wompili.sqlite")
	connection, _ := db.GetConnection("wompili.sqlite")
	err := db.SeedDb(connection, "./db/sql/seed.sql")

	if err != nil {
		log.Fatal("Seeding failed with err: ", err)
	}

	e.GET("/snippets", func(c echo.Context) error {
		snippets, _ := db.GetAllSnippets(connection)
		return c.JSON(http.StatusOK, snippets)
	})

	err = e.Start(":1323")

	if err != nil {
		log.Fatal(err)
	}
}
