package main

import (
	"app/gen/sqlc"
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func initEcho() *echo.Echo {
	e := echo.New()
	e.Debug = true
	e.Logger.SetOutput(os.Stdout)
	return e
}

func newPgxPool(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

type Controller struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func (controller *Controller) GetAuthor(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Logger().Error("parseint: ", err)
		return c.String(http.StatusBadRequest, "parseint: "+err.Error())
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
		return c.String(http.StatusInternalServerError, "list author: "+err.Error())
	}
	return c.JSON(http.StatusOK, authors)
}

func (controller *Controller) InsertAuthor(c echo.Context) error {
	var author sqlc.CreateAuthorParams
	if err := c.Bind(&author); err != nil {
		c.Logger().Error("bind: ", err)
		return c.String(http.StatusBadRequest, "bind: "+err.Error())
	}

	newAuthor, err := controller.queries.CreateAuthor(context.Background(), author)
	if err != nil {
		c.Logger().Error("insert: ", err)
		return c.String(http.StatusBadRequest, "save: "+err.Error())
	}
	c.Logger().Infof("inserted author: %v", newAuthor.ID)
	return c.NoContent(http.StatusCreated)
}

func (controller *Controller) UpdateAuthor(c echo.Context) error {
	// NOTE: このコンテキストをそのまま使って良いのかわかってない
	// コントローラーの他のメソッドでもこうした方が良いのかも
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Logger().Error("parseint: ", err)
		return c.String(http.StatusBadRequest, "parseint: "+err.Error())
	}

	// TODO: pgx.TxOptions の設定について理解したら BeginTx に設定を追加する
	// 以下は Pool.Begin と同じ定義
	c.Logger().Info("begin tx")
	tx, err := controller.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		c.Logger().Error("begin: ", err)
		return c.String(http.StatusInternalServerError, "sorry")
	}
	defer func() {
		if err != nil {
			// err があった場合はロールバックする
			tx.Rollback(ctx)
		} else {
			// err がない場合はコミットする
			tx.Commit(ctx)
		}
	}()

	// トランザクション用の query を取得
	txQuery := controller.queries.WithTx(tx)

	// ロックする
	c.Logger().Info("lock target row")
	author, err := txQuery.LockAuthor(ctx, id)
	if err != nil {
		c.Logger().Error("lock author: ", err)
		return c.String(http.StatusBadRequest, "lock author: "+err.Error())
	}

	// author の値で updateParams を作成
	// NOTE: これにより、リクエストボディに含まれない値は author の値がそのまま使われる
	updateParams := sqlc.UpdateAuthorParams{
		ID:   author.ID,
		Name: author.Name,
		Age:  author.Age,
		Bio:  author.Bio,
	}
	if err := c.Bind(&updateParams); err != nil {
		c.Logger().Error("bind: ", err)
		return c.String(http.StatusBadRequest, "bind: "+err.Error())
	}

	c.Logger().Info("update author")
	newAuthor, err := txQuery.UpdateAuthor(context.Background(), updateParams)
	if err != nil {
		c.Logger().Error("update: ", err)
		return c.String(http.StatusBadRequest, "update: "+err.Error())
	}
	c.Logger().Infof("updated author: %v", newAuthor.ID)
	return c.NoContent(http.StatusOK)
}

func main() {
	pool, err := newPgxPool(os.Getenv("DNS"))
	if err != nil {
		log.Fatal(err)
	}
	queries := sqlc.New(pool)

	controller := &Controller{
		pool:    pool,
		queries: queries,
	}

	e := initEcho()

	e.GET("/api/authors/:id", controller.GetAuthor)
	e.PUT("/api/authors/:id", controller.UpdateAuthor)
	e.GET("/api/authors", controller.ListAuthor)
	e.POST("/api/authors", controller.InsertAuthor)
	e.Logger.Fatal(e.Start(":8989"))
}
