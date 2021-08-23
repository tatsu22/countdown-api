package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/tatsu22/docker-gs-ping/src/node"
	"github.com/tatsu22/docker-gs-ping/src/utils"
)

type DetailedRequest struct {
	Nums []int `json:"nums"`
	Goal int   `json:"goal"`
}

type PlayRequest struct {
	NumSmall int `json:"numSmall"`
	NumBig   int `json:"numBig"`
}

type CompletedGame struct {
	Nums          []int    `json:"nums"`
	Goal          int      `json:"goal"`
	EquationArray []string `json:"equationArray"`
	Equation      string   `json:"equation"`
	Complete      bool     `json:"complete"`
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

	e.POST("/game", func(c echo.Context) error {
		req := new(DetailedRequest)
		if err := c.Bind(req); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, playGameReq(*req))
	})

	e.GET("/game", func(c echo.Context) error {
		return c.JSON(http.StatusOK, playGameReq(genRandomGame()))
	})

	e.GET("/exit", func(c echo.Context) error {
		fmt.Println("Exiting by request")
		os.Exit(0)
		return c.HTML(http.StatusOK, "Exiting")
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

// func simGame(nums []int, goal int) *node.Node {

// 	base := node.GenBaseNode(nums)

// 	base.GenChildren()
// 	return base
// }

func playGameReq(req DetailedRequest) CompletedGame {
	return playGame(req.Goal, req.Nums)
}

func playGame(goal int, playNums []int) CompletedGame {
	logrus.Info("Playing game :", playNums, goal)
	start := time.Now()

	available := []node.Node{}
	base := node.GenBaseNode(playNums)
	base.GenChildren()

	available = append(available, base.Children...)
	calc := node.Node{}
	numNodesCalculated := 0

	for {
		if len(available) == 0 {
			logrus.Warn("No solution!")
			return CompletedGame{}
		}
		calc, available = available[0], available[1:]
		logrus.Debug("Calculating node: ", calc)
		numNodesCalculated++
		calc.GenChildren()
		for _, child := range calc.Children {
			if child.Result == goal {
				logrus.Info("Found result!: ", child)
				elapsed := time.Since(start)
				logrus.Info("Took ", elapsed)
				logrus.Info("Num nodes calculated: ", numNodesCalculated)
				game := new(CompletedGame)
				game.Goal = goal
				game.EquationArray = child.Equation
				game.Nums = playNums
				game.Complete = true
				game.Equation = node.EquationToString(game.EquationArray)
				logrus.Info("Returning game: ", game)
				return *game
			}
			if !node.ContainsNode(available, child) {
				available = node.InsertSorted(available, child, goal)
			}
		}
		if numNodesCalculated%100 == 0 {
			if time.Since(start).Minutes() >= 1.0 {
				logrus.Info("Game could not be completed")
				game := new(CompletedGame)
				game.Nums = playNums
				game.Complete = false
				game.Goal = goal
				return *game
			}
		}
	}
}

func genRandomGame() DetailedRequest {
	availableNums := []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 25, 50, 75, 100}
	numsList := []int{}

	for i := 0; i < 6; i++ {
		index := rand.Intn(len(availableNums))
		numsList = append(numsList, availableNums[index])
		availableNums = utils.RemoveElement(availableNums, index)
	}

	goal := rand.Intn(900) + 100

	return DetailedRequest{Goal: goal, Nums: numsList}
}

func genGame(req PlayRequest) (DetailedRequest, error) {
	if req.NumBig+req.NumSmall != 6 || req.NumBig > 4 || req.NumBig < 0 || req.NumSmall < 2 {
		return DetailedRequest{}, errors.New("Invalid configuration for game")
	}

	smallNumList := []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10}
	bigNumList := []int{25, 50, 75, 100}

	numsList := []int{}

	for i := 0; i < req.NumSmall; i++ {
		index := rand.Intn(len(smallNumList))
		numsList = append(numsList, smallNumList[index])
		smallNumList = utils.RemoveElement(smallNumList, index)
	}

	for i := 0; i < req.NumBig; i++ {
		index := rand.Intn(len(bigNumList))
		numsList = append(numsList, bigNumList[index])
		bigNumList = utils.RemoveElement(bigNumList, index)
	}

	goal := rand.Intn(900) + 100

	return DetailedRequest{Goal: goal, Nums: numsList}, nil
}
