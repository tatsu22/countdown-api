package main

import (
	"fmt"
	"math/rand"
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
	Nums []int `json:"nums"`
	Goal int   `json:"goal"`
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
	Message string `json:"message"`
	Error   error  `json:"errorMsg"`
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
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "'nums' query parameter missing",
			})
		}
		nums := make([]int, len(numsStrings))
		var err error
		for i, v := range numsStrings {
			nums[i], err = strconv.Atoi(v)
			if err != nil {
				logrus.Info("Error when converting nums string to ints", err)
				return c.JSON(http.StatusBadRequest, ErrorResponse{
					Message: "Could not convert nums to integers",
					Error:   err,
				})
			}
		}

		goal, err := strconv.Atoi(c.QueryParam("goal"))
		if err != nil {
			logrus.Info("Error when converting goal to ints", err)
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "Could not convert goal to integer",
				Error:   err,
			})
		}

		// shortestString := c.QueryParam("shortest")
		// var shortest bool
		// if len(shortestString) == 0 {
		// 	logrus.Debug("Defaulting shortest to false")
		// 	shortest = false
		// } else {
		// 	shortest, err = strconv.ParseBool(shortestString)
		// 	if err != nil {
		// 		logrus.Info("Error when converting shortest to bool", err)
		// 		return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "Could not parse shortest to bool"})
		// 	}
		// }

		req := DetailedRequest{Nums: nums, Goal: goal}
		return c.JSON(http.StatusOK, playGameReq(req))
	})

	e.GET("/createGame", func(c echo.Context) error {
		numSmall, err := strconv.Atoi(c.QueryParam("numSmall"))
		if err != nil {
			logrus.Info("Error when parsing numSmall query param: ", c.QueryParam("numSmall"))
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "Could not convert numSmall to integer",
				Error:   err,
			})
		}

		numBig, err := strconv.Atoi(c.QueryParam("numBig"))
		if err != nil {
			logrus.Info("Error when parsing numBig query param: ", c.QueryParam("numBig"))
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "Could not convert numBig to integer",
				Error:   err,
			})
		}

		game, err := genGame(numSmall, numBig)
		if err != nil {
			logrus.Info("Error when generating game based on small and big nums: ", numSmall, numBig)
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "Could not generate game based on input",
				Error:   err,
			})
		}
		return c.JSON(http.StatusOK, game)
	})

	httpPort := os.Getenv("COUNTDOWN_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func playGameReq(req DetailedRequest) CompletedGame {
	return playGame(req.Goal, req.Nums)
}

func playGame(goal int, playNums []int) CompletedGame {
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
			return CompletedGame{
				Nums:            playNums,
				Goal:            goal,
				Complete:        false,
				TimeTaken:       time.Since(start).String(),
				NodesCalculated: numNodesCalculated,
			}
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
			// Checking for duplicates and sorting takes long enough that it doesn't seem worth it
			// For reference, with below if statement code it seems to take ~1 second to calculate 1000
			// nodes(probably longer for more difficult games), and without the if statement it takes
			// about 10ms to calculate 1000 nodes, and it seems to keep that speed even for solves that
			// that take about 500k nodes. If I find a use case where it's better to sort,
			// then I will uncomment this code
			// if !ContainsNode(available, child) {
			// available = InsertSorted(available, child, goal, shortest)
			// }
		}
		available = append(available, newNodes...)
		if numNodesCalculated%10000 == 0 {
			if time.Since(start).Minutes() >= 1.0 {
				logrus.Info("Game could not be completed")
				return CompletedGame{
					Nums:            playNums,
					Complete:        false,
					Goal:            goal,
					TimeTaken:       time.Since(start).String(),
					NodesCalculated: numNodesCalculated,
				}
			}
		}
	}
}

func genGame(numSmall, numBig int) (DetailedRequest, error) {
	if numBig+numSmall != 6 || numBig > 4 || numBig < 0 || numSmall < 2 {
		return DetailedRequest{}, GenGameError{
			Msg:      "Could not generate game based on inputs",
			NumSmall: numSmall,
			NumBig:   numBig,
		}
	}

	smallNumList := []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10}
	bigNumList := []int{25, 50, 75, 100}

	numsList := []int{}

	for i := 0; i < numSmall; i++ {
		index := rand.Intn(len(smallNumList))
		numsList = append(numsList, smallNumList[index])
		smallNumList = RemoveElement(smallNumList, index)
	}

	for i := 0; i < numBig; i++ {
		index := rand.Intn(len(bigNumList))
		numsList = append(numsList, bigNumList[index])
		bigNumList = RemoveElement(bigNumList, index)
	}

	goal := rand.Intn(900) + 100

	return DetailedRequest{Goal: goal, Nums: numsList}, nil
}

type GenGameError struct {
	Msg      string `json:"message"`
	NumSmall int    `json:"numSmall"`
	NumBig   int    `json:"numBig"`
}

func (g GenGameError) Error() string {
	return fmt.Sprint("Could not generate game where smallNums=%s and bigNums=%s", g.NumSmall, g.NumBig)
}
