package main

import (
	"fmt"

	"github.com/tatsu22/docker-gs-ping/src/node"
)

func main2() {
	fmt.Println(node.EquationToString([]string{"3", "4", "*", "1", "+"}))
}
