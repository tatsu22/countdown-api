package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/olliefr/docker-gs-ping/src/node"
)

type Request struct {
	nums      []int
	smallNums int
	largeNums int
	target    int
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.GET("/game", func(c echo.Context) error {
		return c.HTML(http.StatusOK, simGame([]int{1, 2, 3, 4, 10, 25}, 156))
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func simGame(nums []int, goal int) string {

	base := node.GenBaseNode(nums)

	for {
		base.GenChildren()
		return fmt.Sprintf("Num children %d", len(base.GetChildren()))
	}

	return ""
}
