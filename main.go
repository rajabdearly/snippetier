package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"snippetier/db"
	"snippetier/routes"
)

const dbName = "wompili.sqlite"

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	e.GET("/", routes.Hello)

	db.SetupNewTestDb(dbName)
	storage, err := db.GetConnection(dbName)
	defer storage.CloseConnection()

	if err != nil {
		log.Fatal("Seeding failed with err: ", err)
	}

	e.GET("/snippets", routes.GetAllSnippets(storage))

	e.POST("/snippets/new", routes.SaveSnippet(storage))
	e.PUT("/snippets/:id", routes.UpdateSnippet(storage))
	e.DELETE("/snippets/:id", routes.DeleteSnippet(storage))

	err = e.Start(":1323")

	if err != nil {
		log.Fatal(err)
	}
}
