package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"task/app/api"
	"time"
)

// Calculate - функция для вычисления результата выражения
func Calculate(exp string) (int, error) {
	postfix := InfixToPostfix(exp)
	if postfix == "" {
		return 0, fmt.Errorf("полученное постфиксное выражение пусто")
	}
	result := EvaluatePostfix(postfix)
	return result, nil
}

// GetOperations получает список операций и их длительности
func GetOperations() ([]api.OperationDuration, error) {
	url := "http://localhost:5500/operations/get"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении HTTP-запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("получен неправильный статус код: %d", resp.StatusCode)
	}

	var operations []api.OperationDuration
	err = json.NewDecoder(resp.Body).Decode(&operations)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования JSON: %v", err)
	}

	return operations, nil
}

// InfixToPostfix - функция для преобразования инфиксного выражения в постфиксное
func InfixToPostfix(expression string) string {
	tokens := strings.Fields(expression)
	outputQueue := []string{}
	operatorStack := []string{}

	for _, token := range tokens {
		if isNumber(token) {
			outputQueue = append(outputQueue, token)
		} else if isOperator(token) {
			for len(operatorStack) > 0 && precedence(operatorStack[len(operatorStack)-1]) >= precedence(token) {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, token)
		}
	}

	for len(operatorStack) > 0 {
		outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return strings.Join(outputQueue, " ")
}

// EvaluatePostfix - функция для вычисления результата постфиксного выражения
func EvaluatePostfix(expression string) int {
	tokens := strings.Fields(expression)
	stack := []int{}

	for _, token := range tokens {
		if isNumber(token) {
			num, _ := strconv.Atoi(token)
			stack = append(stack, num)
		} else if isOperator(token) {
			operand2 := stack[len(stack)-1]
			operand1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			durations, err := GetOperations()
			if err != nil {
				log.Printf("Error fetching operations: %v", err)
				return 0
			}

			result := applyOperator(token, operand1, operand2, durations)
			stack = append(stack, result)
		}
	}

	return stack[0]
}

// isNumber - функция для проверки, является ли строка числом
func isNumber(token string) bool {
	_, err := strconv.Atoi(token)
	return err == nil
}

// isOperator - функция для проверки, является ли строка оператором
func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

// precedence - функция для определения приоритета оператора
func precedence(operator string) int {
	switch operator {
	case "*", "/":
		return 2
	case "+", "-":
		return 1
	default:
		return 0
	}
}

// getDuration - функция для получения длительности операции
func getDuration(sign string, durations []api.OperationDuration) int {
	for _, op := range durations {
		if op.OperationName == sign {
			return op.OperationDuration
		}
	}
	return 0
}

func applyOperator(operator string, operand1, operand2 int, durations []api.OperationDuration) int {
	switch operator {
	case "+":
		timeMy := getDuration("plus", durations)
		timerDuration := time.Duration(timeMy) * time.Millisecond
		timer := time.NewTimer(timerDuration)
		<-timer.C
		return operand1 + operand2

	case "-":
		timeMy := getDuration("minus", durations)
		timerDuration := time.Duration(timeMy) * time.Millisecond
		timer := time.NewTimer(timerDuration)
		<-timer.C
		return operand1 - operand2

	case "*":
		timeMy := getDuration("multiply", durations)
		timerDuration := time.Duration(timeMy) * time.Millisecond
		timer := time.NewTimer(timerDuration)
		<-timer.C
		return operand1 * operand2

	case "/":
		timeMy := getDuration("divide", durations)
		timerDuration := time.Duration(timeMy) * time.Millisecond
		timer := time.NewTimer(timerDuration)
		<-timer.C
		return operand1 / operand2

	default:
		return 0
	}
}
