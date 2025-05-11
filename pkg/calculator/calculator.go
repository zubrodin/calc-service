package calculator

import (
	"fmt"
	"strconv"
	"strings"
)

type Calculator struct{}

func New() *Calculator {
	return &Calculator{}
}

func (c *Calculator) Calculate(expr string) (float64, error) {
	expr = strings.ReplaceAll(expr, " ", "")

	rpn, err := c.toRPN(expr)
	if err != nil {
		return 0, err
	}

	result, err := c.evaluateRPN(rpn)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (c *Calculator) toRPN(expr string) ([]string, error) {
	var output []string
	var operators []string

	for i := 0; i < len(expr); {
		char := string(expr[i])

		if c.isDigit(char) {
			j := i
			for j < len(expr) && (c.isDigit(string(expr[j])) || string(expr[j]) == ".") {
				j++
			}
			output = append(output, expr[i:j])
			i = j
		} else if char == "(" {
			operators = append(operators, char)
			i++
		} else if char == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			operators = operators[:len(operators)-1]
			i++
		} else if c.isOperator(char) {
			for len(operators) > 0 && c.precedence(operators[len(operators)-1]) >= c.precedence(char) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, char)
			i++
		} else {
			return nil, fmt.Errorf("invalid character: %s", char)
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

func (c *Calculator) evaluateRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		if c.isNumber(token) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("invalid expression")
			}
			a := stack[len(stack)-2]
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				result = a / b
			default:
				return 0, fmt.Errorf("unknown operator: %s", token)
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	return stack[0], nil
}

func (c *Calculator) isDigit(s string) bool {
	return s >= "0" && s <= "9"
}

func (c *Calculator) isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (c *Calculator) isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func (c *Calculator) precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}
