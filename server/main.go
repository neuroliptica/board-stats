package main

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	db Database

	AvaibleBoards = map[string]bool{
		"2ch.hk/b":  true,
		"2ch.hk/po": true,
		"2ch.hk/vg": true,
	}
)

func init() {
	b, err := sql.Open("sqlite3", "stats.db")
	if err != nil {
		panic(err)
	}
	db = Database{b}
	err = db.MigrateStats()
	if err != nil {
		panic(err)
	}

	sosach := Daemon{
		f: func() (*Record, error) {
			return FetchSosach("b")
		},
		table: "2ch.hk/b",
	}
	go sosach.Run()
}

func logging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Info().Fields(map[string]interface{}{
			"method": c.Request().Method,
			"url":    c.Request().URL.Path,
			"query":  c.Request().URL.RawQuery,
		}).Msg("request")

		err := next(c)
		if err != nil {
			logger.Error().Fields(map[string]interface{}{
				"error": err.Error(),
			}).Msg("response")
			return err
		}
		return nil
	}
}

func main() {
	e := echo.New()
	logger = NewLogger()
	e.Use(logging, middleware.Recover(), middleware.CORS())

	e.Static("/static", "static")

	e.GET("/api/stats", func(c echo.Context) error {
		type FailedResponse struct {
			Error string `json:"error"`
		}
		board := c.QueryParam("board")
		if ok := AvaibleBoards[board]; !ok {
			return c.JSON(http.StatusNotFound, FailedResponse{
				Error: "no such board to track",
			})
		}
		res, err := db.GetBoardStats(board)
		if err != nil {
			return c.JSON(http.StatusNotFound, FailedResponse{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, res)
	})

	logger.Info().Msg(e.Start(":8080").Error())
}
