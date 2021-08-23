package main

import (
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
)

var Operations = [4]string{"+", "-", "*", "/"}

type Node struct {
	Result        int `json:"result"`
	Equation      []string
	Children      []Node
	RemainingNums []int
}

func (node *Node) GenChildren() {
	logrus.Debug("Generating children for node: ", node)
	logrus.Debug("Going through possible operations")

	// Should be able to replace this by just checking length of equation == 11
	// if len(node.Equation) >= 1 {
	// 	_, err := utils.IsNumber(node.Equation[len(node.Equation)-1])
	// 	if len(node.RemainingNums) == 0 && err != nil {
	// 		logrus.Debug("Generating no new children, node is at tip")
	// 		return
	// 	}
	// }

	// Replaced above by just checking length of a full equation, seems to work
	if len(node.Equation) >= 11 {
		logrus.Debug("Generating no new children, node is at tip")
		logrus.Debug("Node: ", node)
		return
	}

	for _, v := range Operations {
		if IsValidEquation(append(node.Equation, v)) {
			newEquation := make([]string, len(node.Equation))
			copy(newEquation, node.Equation)
			newEquation = append(newEquation, v)
			logrus.Debug("Equation is valid: ", newEquation)
			if len(newEquation) > 2 {
				if IsCompleteEquation(newEquation) {
					newResult, err := EvaluateEquation(newEquation)
					if err == nil {
						logrus.Debug("Equation is valid and complete, result is: ", newResult)
						node.Children = append(node.Children, Node{Equation: newEquation, RemainingNums: node.RemainingNums, Result: newResult})
					}
				} else {
					logrus.Debug("Equation is valid but incomplete, setting result to parent: ", node.Result)
					node.Children = append(node.Children, Node{Equation: newEquation, RemainingNums: node.RemainingNums, Result: node.Result})
				}
			} else {
				logrus.Warn("You should not see this, added operation and length <= 2")
				node.Children = append(node.Children, Node{Equation: newEquation, RemainingNums: node.RemainingNums, Result: 0})
			}
			// logrus.Debug("added node: ", node.Children[len(node.Children)-1])
		}
	}

	// For numbers there is much less garbage to do, since we know
	// 1. they are valid
	// 2. they are are not complete, and
	// 3. we don't care about the length if we're adding a number
	logrus.Debug("Going through possible nums: ", node.RemainingNums)
	for i, v := range node.RemainingNums {
		newEquation := make([]string, len(node.Equation))
		copy(newEquation, node.Equation)
		newEquation = append(newEquation, fmt.Sprint(v))
		newRemainingNums := RemoveElement(node.RemainingNums, i)
		node.Children = append(node.Children, Node{Equation: newEquation, RemainingNums: newRemainingNums, Result: node.Result})
	}

	logrus.Debug("Generated children for node: ", node)
	logrus.Debug("Children: ")
	for _, v := range node.Children {
		logrus.Debug(v)
	}
}

// Generates node with no equation and all remaining nums
func GenBaseNode(nums []int) *Node {
	return &Node{Equation: []string{}, RemainingNums: nums, Result: 0}
}

func InsertSorted(s []Node, e Node, goal int) []Node {
	// This will not guarantee shortest solution, but is faster than looking for shortest
	i := sort.Search(len(s), func(i int) bool { return AbsVal(s[i].Result-goal) > AbsVal(e.Result-goal) })

	// This will give us shortest solution guaranteed
	// i := sort.Search(len(s), func(i int) bool {
	// 	if len(s[i].Equation) == len(e.Equation) {
	// 		return utils.AbsVal(s[i].Result-goal) > utils.AbsVal(e.Result-goal)
	// 	}
	// 	return len(s[i].Equation) > len(e.Equation)
	// })

	// If we get back the length then it was unsorted, have to append b/c insertion logic below will break if we give it the length
	if i == len(s) {
		s = append(s, e)
		return s
	}
	s = append(s, Node{})
	copy(s[i+1:], s[i:])
	s[i] = e
	return s
}

func Compare(a, b Node) bool {
	if len(a.RemainingNums) != len(b.RemainingNums) || len(a.Equation) != len(b.Equation) || a.Result != b.Result {
		return false
	}

	// Here we know all lists are same length and result is the same

	// Check for equation equality
	for i := 0; i < len(a.Equation); i++ {
		if a.Equation[i] != b.Equation[i] {
			return false
		}
	}

	// Check for remaining nums equality
	for i := 0; i < len(a.RemainingNums); i++ {
		if a.RemainingNums[i] != b.RemainingNums[i] {
			return false
		}
	}

	return true
}

func ContainsNode(list []Node, a Node) bool {
	for _, v := range list {
		if Compare(a, v) {
			return true
		}
	}
	return false
}

func EquationToString(equation []string) string {
	queue := []string{}

	for _, v := range equation {
		_, err := IsNumber(v)
		if err == nil {
			queue = append(queue, v)
		} else {
			length := len(queue)
			newInfix := "(" + queue[length-1] + v + queue[length-2] + ")"
			queue = append(queue[:length-2], newInfix)
		}
	}
	return queue[0]
}
