package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type DetailedRequest struct {
	Nums     []int `json:"nums"`
	Goal     int   `json:"goal"`
	Shortest bool  `json:"shortest"`
}

type PlayRequest struct {
	NumSmall int `json:"numSmall"`
	NumBig   int `json:"numBig"`
}

type CompletedGame struct {
	Nums            []int    `json:"nums"`
	Goal            int      `json:"goal"`
	EquationArray   []string `json:"equationArray"`
	Equation        string   `json:"equation"`
	Complete        bool     `json:"complete"`
	TimeTaken       string   `json:"timeTaken"`
	NodesCalculated int      `json:"nodesCalculated"`
}

type ErrorResponse struct {
	Error string `json:"errorMsg"`
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

	// TODO: Use query parameters to generate game
	e.GET("/game", func(c echo.Context) error {
		numsStrings := strings.Split(c.QueryParam("nums"), ",")
		if len(numsStrings) != 6 {
			logrus.Info("Missing nums query parameter")
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "'nums' query parameter missing"})
		}
		nums := make([]int, len(numsStrings))
		var err error
		for i, v := range numsStrings {
			nums[i], err = strconv.Atoi(v)
			if err != nil {
				logrus.Info("Error when converting nums string to ints", err)
				return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "Could not convert nums to integers"})
			}
		}

		goal, err := strconv.Atoi(c.QueryParam("goal"))
		if err != nil {
			logrus.Info("Error when converting goal to ints", err)
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "Could not convert goal to integer"})
		}

		shortestString := c.QueryParam("shortest")
		var shortest bool
		if len(shortestString) == 0 {
			logrus.Debug("Defaulting shortest to false")
			shortest = false
		} else {
			shortest, err = strconv.ParseBool(shortestString)
			if err != nil {
				logrus.Info("Error when converting shortest to bool", err)
				return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "Could not parse shortest to bool"})
			}
		}

		req := DetailedRequest{Nums: nums, Goal: goal, Shortest: shortest}
		return c.JSON(http.StatusOK, playGameReq(req))
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func playGameReq(req DetailedRequest) CompletedGame {
	return playGame(req.Goal, req.Nums, req.Shortest)
}

func playGame(goal int, playNums []int, shortest bool) CompletedGame {
	logrus.Info("Playing game :", playNums, goal)
	start := time.Now()

	available := []Node{}
	base := GenBaseNode(playNums)
	baseChildren := base.GenChildren()

	available = append(available, baseChildren...)
	calc := Node{}
	numNodesCalculated := 0

	for {
		if len(available) == 0 {
			logrus.Warn("No solution!")
			return CompletedGame{}
		}
		calc, available = available[0], available[1:]
		logrus.Debug("Calculating node: ", calc)
		numNodesCalculated++
		newNodes := calc.GenChildren()
		for _, child := range newNodes {
			if child.Result == goal {
				logrus.Info("Found result!: ", child)
				elapsed := time.Since(start)
				logrus.Info("Took ", elapsed)
				logrus.Info("Num nodes calculated: ", numNodesCalculated)
				game := CompletedGame{
					Goal:            goal,
					EquationArray:   child.Equation,
					Nums:            playNums,
					Complete:        true,
					Equation:        EquationToString(child.Equation),
					TimeTaken:       elapsed.String(),
					NodesCalculated: numNodesCalculated,
				}
				logrus.Info("Returning game: ", game)
				return game
			}
			if !ContainsNode(available, child) {
				// available = InsertSorted(available, child, goal, shortest)
				available = append(available, child)
			}
		}
		if numNodesCalculated%1000 == 0 {
			if time.Since(start).Minutes() >= 1.0 {
				logrus.Info("Game could not be completed")
				return CompletedGame{
					Nums:     playNums,
					Complete: false,
					Goal:     goal,
				}
			}
		}
	}
}

// func genGame(req PlayRequest) (DetailedRequest, error) {
// 	if req.NumBig+req.NumSmall != 6 || req.NumBig > 4 || req.NumBig < 0 || req.NumSmall < 2 {
// 		return DetailedRequest{}, errors.New("invalid configuration for game")
// 	}

// 	smallNumList := []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10}
// 	bigNumList := []int{25, 50, 75, 100}

// 	numsList := []int{}

// 	for i := 0; i < req.NumSmall; i++ {
// 		index := rand.Intn(len(smallNumList))
// 		numsList = append(numsList, smallNumList[index])
// 		smallNumList = RemoveElement(smallNumList, index)
// 	}

// 	for i := 0; i < req.NumBig; i++ {
// 		index := rand.Intn(len(bigNumList))
// 		numsList = append(numsList, bigNumList[index])
// 		bigNumList = RemoveElement(bigNumList, index)
// 	}

// 	goal := rand.Intn(900) + 100

// 	return DetailedRequest{Goal: goal, Nums: numsList}, nil
// }
