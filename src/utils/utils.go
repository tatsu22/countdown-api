package utils

import (
	"errors"
)

func IsValidEquation(equation []byte) bool {

	if !IsNumber(equation[0]) || !IsNumber(equation[1]) {
		return false
	}

	numNums := 0
	numOps := 0

	for i := 0; i < len(equation); i++ {
		if IsNumber(equation[i]) {
			numNums++
		} else {
			numOps++
			if numOps >= numNums {
				return false
			}
		}
	}

	return true
}

func EvaluateEquation(equation []byte) (int, error) {
	queue := []int{}

	for _, v := range equation {
		if IsNumber(v) {
			queue = append(queue, int(v))
		} else {
			var method = add
			length := len(queue)
			switch v {
			case '*':
				method = multiply
			case '+':
				method = add
			case '-':
				method = sub
			case '/':
				method = divide
			}
			newNum, err := method(queue[length-1], queue[length-2])
			if err != nil {
				return -1, errors.New("equation resulted in error")
			} else {
				queue = append(queue[:length-2], newNum)
			}
		}
	}

	return queue[len(queue)-1], nil
}

func add(a, b int) (int, error) {
	return a + b, nil
}

func sub(a, b int) (int, error) {
	if a-b < 0 {
		return -1, errors.New("subtraction resulted in negative number")
	}
	return a - b, nil
}

func multiply(a, b int) (int, error) {
	return a * b, nil
}

func divide(a, b int) (int, error) {
	if a%b == 0 {
		return a / b, nil
	}
	return -1, errors.New("division resulted in remainder")
}

func IsNumber(q byte) bool {
	return '0' <= q && q <= '9'
}

func RemoveElement(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
