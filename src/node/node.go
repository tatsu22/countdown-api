package node

import (
	"github.com/olliefr/docker-gs-ping/src/utils"
)

var Operations = [4]byte{'+', '-', '*', '/'}

type Node struct {
	result        int
	equation      []byte
	children      []Node
	remainingNums []int
}

func (node *Node) GenChildren() {
	for _, v := range Operations {
		if utils.IsValidEquation(append(node.equation, v)) {
			newEquation := append(node.equation, v)
			newResult, err := utils.EvaluateEquation(newEquation)
			if err == nil {
				node.children = append(node.children, Node{equation: newEquation, remainingNums: node.remainingNums, result: newResult})
			}
		}
	}

	for i, v := range node.remainingNums {
		if utils.IsValidEquation(append(node.equation, byte(v))) {
			newEquation := append(node.equation, byte(v))
			newResult, err := utils.EvaluateEquation(newEquation)
			if err == nil {
				newRemainingNums := utils.RemoveElement(node.remainingNums, i)
				node.children = append(node.children, Node{equation: newEquation, remainingNums: newRemainingNums, result: newResult})
			}
		}
	}
}

func GenBaseNode(nums []int) Node {
	return Node{equation: []byte{}, remainingNums: nums, result: 0}
}

func (node *Node) GetChildren() []Node {
	return node.children
}
