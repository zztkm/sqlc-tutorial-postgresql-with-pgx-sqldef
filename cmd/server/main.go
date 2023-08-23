package main

import (
	"app/gen/sqlc"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func initEcho() *echo.Echo {
	e := echo.New()
	e.Debug = true
	e.Logger.SetOutput(os.Stdout)
	return e
}

type Controller struct {
	queries *sqlc.Queries
}

func (controller *Controller) GetAuthor(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Logger().Error("parseint: ", err)
		return c.String(http.StatusBadRequest, "parseint: " + err.Error())
	}

	author, err := controller.queries.GetAuthor(context.Background(), id)
	if err != nil {
		c.Logger().Error("get author: ", err)
		return c.String(http.StatusInternalServerError, "not found")
	}
	return c.JSON(http.StatusOK, author)
}

func (controller *Controller) ListAuthor(c echo.Context) error {
	authors, err := controller.queries.ListAuthors(context.Background())
	if err != nil {
		c.Logger().Error("list author: ", err)
		return c.String(http.StatusInternalServerError, "list author: " + err.Error())
	}
	return c.JSON(http.StatusOK, authors)
}

func (controller *Controller) InsertAuthor(c echo.Context) error {
	var author sqlc.CreateAuthorParams
	if err := c.Bind(&author); err != nil {
		c.Logger().Error("bind: ", err)
		return c.String(http.StatusBadRequest, "bind: " + err.Error())
	}

	newAuthor, err := controller.queries.CreateAuthor(context.Background(), author)
	if err != nil {
		c.Logger().Error("insert: ", err)
		return c.String(http.StatusBadRequest, "save: " + err.Error())
	}
	c.Logger().Infof("inserted author: %v", newAuthor.ID)
	return c.NoContent(http.StatusCreated)
}

func main() {
	fmt.Println(os.Getenv("DNS"))
	conn, err := pgx.Connect(context.Background(), os.Getenv("DNS"))
	if err != nil {
		log.Fatal(err)
	}
	queries := sqlc.New(conn)

	controller := &Controller{
		queries: queries,
	}

	e := initEcho()

	e.GET("/api/authors/:id", controller.GetAuthor)
	e.GET("/api/authors", controller.ListAuthor)
	e.POST("/api/authors", controller.InsertAuthor)
	e.Logger.Fatal(e.Start(":8989"))
}
