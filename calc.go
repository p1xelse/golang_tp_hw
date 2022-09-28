package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"unicode"
)

const (
	openBracket  byte = '('
	closeBracket byte = ')'
	plus         byte = '+'
	minus        byte = '-'
	multiply     byte = '*'
	divide       byte = '/'
)

func getPriority(char byte) int {
	priorityMap := map[byte]int{
		'(': 0,
		')': 0,
		'+': 1,
		'-': 1,
		'/': 2,
		'*': 2,
	}

	return priorityMap[char]
}

func isOperator(char byte) bool {
	return char == openBracket ||
		char == closeBracket ||
		char == plus ||
		char == minus ||
		char == multiply ||
		char == divide
}

func readOperand(reader *strings.Reader) (f float64, err error) {
	offset := reader.Size() - int64(reader.Len())
	_, err = fmt.Fscan(reader, &f)

	if err != nil {
		reader.Seek(offset, io.SeekStart)
	}

	return
}

func skipSpaces(reader *strings.Reader) error {
	for {
		b, err := reader.ReadByte()

		if err != nil {
			return err
		}

		if !unicode.IsSpace(rune(b)) {
			_ = reader.UnreadByte()
			return nil
		}
	}
}

func readOperator(reader *strings.Reader) (char byte, err error) {
	err = skipSpaces(reader)
	if err != nil {
		return
	}

	char, err = reader.ReadByte()

	if err != nil {
		return
	}

	if !isOperator(char) {
		err = errors.New("invalid expression: operator expected")
	}

	return
}

func infixToPostfix(infixExpr string) (postfixExpr string, err error) {
	reader := strings.NewReader(infixExpr)
	var operators Stack

	checkOperand := true

	for reader.Len() != 0 {
		if checkOperand {
			operand, err := readOperand(reader)
			if err == nil {
				postfixExpr = fmt.Sprintf("%s%f ", postfixExpr, operand)
				checkOperand = false
				continue
			}
		}

		operator, err := readOperator(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		checkOperand = true

		if operator == openBracket {
			operators.Push(operator)
			continue
		} else if operator == closeBracket {
			checkOperand = false

			for {
				if operators.isEmpty() {
					return "", errors.New("invalid expression")
				}
				operator = operators.Pop().(byte)
				if operator == openBracket {
					break
				}
				postfixExpr = fmt.Sprintf("%s%s ", postfixExpr, string(operator))
			}
			continue
		}

		if !operators.isEmpty() && getPriority(operator) <= getPriority(operators.Top().(byte)) {
			postfixExpr = fmt.Sprintf("%s%s ", postfixExpr, string(operators.Pop().(byte)))
		}

		operators.Push(operator)
	}

	for !operators.isEmpty() {
		operator := operators.Pop().(byte)

		if operator == openBracket {
			return "", errors.New("invalid expression")
		}
		postfixExpr = fmt.Sprintf("%s%s ", postfixExpr, string(operator))
	}

	postfixExpr = strings.TrimSpace(postfixExpr)

	return
}

func apply(operands *Stack, operator byte) error {
	if operands.isEmpty() {
		return errors.New("invalid format")
	}
	b := operands.Pop().(float64)
	if operands.isEmpty() {
		return errors.New("invalid format")
	}
	a := operands.Pop().(float64)
	c, err := execute(a, b, operator)
	if err != nil {
		return err
	}
	operands.Push(c)
	return nil
}

func execute(a float64, b float64, operator byte) (float64, error) {
	switch operator {
	case plus:
		return a + b, nil
	case minus:
		return a - b, nil
	case multiply:
		return a * b, nil
	case divide:
		if math.Abs(b) < 1e-7 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	}
	return 0, errors.New("invalid operand") // never
}

func Calc(str string) (result float64, err error) {
	postfixExpr, err := infixToPostfix(str)

	if err != nil {
		return 0, err
	}

	reader := strings.NewReader(postfixExpr)

	var operands Stack

	for reader.Len() != 0 {
		operand, err := readOperand(reader)

		var operator byte

		if err == nil {
			operands.Push(operand)
			continue
		} else {
			operator, _ = readOperator(reader)
		}
		err = apply(&operands, operator)

		if err != nil {
			return 0, err
		}
	}

	if operands.isEmpty() {
		return 0, errors.New("empty expression")
	}
	result = operands.Pop().(float64)

	if !operands.isEmpty() {
		return 0, errors.New("invalid expression")
	}
	return
}
