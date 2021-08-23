package utils

import (
	"errors"
	"strconv"

	"github.com/sirupsen/logrus"
)

func IsValidEquation(equation []string) bool {

	if len(equation) == 1 {
		_, err := IsNumber(equation[0])
		// logrus.Info("Length is 1, and equation is valid: ", err == nil)
		return err == nil
	}

	if len(equation) == 2 {
		_, err1 := IsNumber(equation[0])
		_, err2 := IsNumber(equation[1])

		// logrus.Info("equation is length 2, both are numbers: ", err1 == nil && err2 == nil)
		return err1 == nil && err2 == nil
	}

	numNums := 0
	numOps := 0

	for i := 0; i < len(equation); i++ {
		_, err := IsNumber(equation[i])
		if err == nil {
			numNums++
		} else {
			numOps++
			if numOps >= numNums {
				// logrus.Info("Number of operation bits is equal or greater than number of numbers")
				return false
			}
		}
	}

	return true
}

func IsCompleteEquation(equation []string) bool {
	numNums := 0
	numOps := 0
	for i := 0; i < len(equation); i++ {
		_, err := IsNumber(equation[i])
		if err == nil {
			numNums++
		} else {
			numOps++
			if numOps >= numNums {
				// logrus.Info("Number of operation bits is equal or greater than number of numbers")
				return false
			}
		}
	}

	if numNums == numOps+1 {
		logrus.Debug("Equation is complete")
		return true
	} else {
		logrus.Debug("Equation is incomplete")
		return false
	}

}

func EvaluateEquation(equation []string) (int, error) {
	queue := []int{}

	for _, v := range equation {
		newNum, err := IsNumber(v)

		if err == nil {
			queue = append(queue, newNum)
		} else {
			var method = add
			length := len(queue)
			if length < 2 {
				return -1, errors.New("not enough numbers for operation")
			}
			switch v {
			case "*":
				method = multiply
			case "+":
				method = add
			case "-":
				method = sub
			case "/":
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
	if a-b <= 0 {
		return -1, errors.New("subtraction resulted in 0 or negative number")
	}
	return a - b, nil
}

func multiply(a, b int) (int, error) {
	if a == 1 || b == 1 {
		return -1, errors.New("multiplying by 1 is useless")
	}
	return a * b, nil
}

func divide(a, b int) (int, error) {
	if a%b == 0 && b != 0 && a != b && a != 1 && b != 1 {
		return a / b, nil
	}
	return -1, errors.New("division resulted in remainder, dividing by 0 or 1, or dividing by self")
}

func IsNumber(q string) (int, error) {
	return strconv.Atoi(q)
}

func RemoveElement(s []int, i int) []int {
	newArray := make([]int, len(s))
	copy(newArray, s)
	newArray[i] = newArray[len(newArray)-1]
	return newArray[:len(newArray)-1]
}

func AbsVal(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
