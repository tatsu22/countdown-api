package main

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

const (
	connectionString = "postgres://postgres:secret@192.168.0.1:5433/countdown_api?sslmode=disable"
)

func Insert(game CompletedGame) bool {
	conn, _ := pgx.Connect(context.Background(), connectionString)

	defer conn.Close(context.Background())

	sqlStatement := `
INSERT INTO completed_game (nums, goal, equation_array, equation, complete, time_taken, nodes_calculated)
VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := conn.Exec(context.Background(), sqlStatement, game.Nums, game.Goal, game.EquationArray, game.Equation, game.Complete, game.TimeTaken, game.NodesCalculated)

	if err != nil {
		panic(err)
	}
	return true
}

func GetGame(nums []int, goal int) CompletedGame {
	conn, _ := pgx.Connect(context.Background(), connectionString)
	defer conn.Close(context.Background())

	sqlStatement := `SELECT nums, goal, equation_array, equation, complete, time_taken, nodes_calculated FROM completed_game WHERE nums=$1 AND goal=$2`

	var completedGame CompletedGame

	// if err := conn.QueryRow(context.Background(), sqlStatement, nums, goal).Scan(&completedGame); err != nil {

	if err := conn.QueryRow(context.Background(), sqlStatement, nums, goal).Scan(&completedGame.Nums,
		&completedGame.Goal,
		&completedGame.EquationArray,
		&completedGame.Equation,
		&completedGame.Complete,
		&completedGame.TimeTaken,
		&completedGame.NodesCalculated); err != nil {
		logrus.Error("Error occured while retrieving game: ", err)
	}
	return completedGame
}

func GetAllGames() []CompletedGame {
	conn, _ := pgx.Connect(context.Background(), connectionString)
	// if err != nil {
	// 	panic(err)
	// }
	defer conn.Close(context.Background())

	var returnList []CompletedGame

	if rows, err := conn.Query(context.Background(), "SELECT * FROM completed_game"); err != nil {
		logrus.Error("Unable to connect to DB", err)
		return []CompletedGame{}
	} else {
		defer rows.Close()

		var tmp CompletedGame
		for rows.Next() {
			rows.Scan(&tmp.Nums, &tmp.Goal, &tmp.EquationArray, &tmp.Equation, &tmp.Complete, &tmp.TimeTaken, &tmp.NodesCalculated)
			returnList = append(returnList, tmp)
		}

		if rows.Err() != nil {
			logrus.Error("Error reading table: ", err)
		}
	}

	return returnList
}
