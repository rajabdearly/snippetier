package main

import (
	"fmt"
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

	routes.SetupRoutes(e, storage, config)
	fmt.Println(config)

	e.GET("/", rootHandler)

	err = e.Start(":1323")
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, braaaat!")
}
