package api

import (
	"fmt"
	"time"
)

type OperationDuration struct {
	OperationName     string // plus minus multiply divide agent
	OperationDuration int    // Длительность операции
}

type Response struct {
	Status int `json:"status"` // Статус http
	Resp   any `json:"result"` // Ответ, содержащий зачастую бесполезную информацию
}

type Expression struct {
	Status         string // Status: Active, Done, Bad,ServerError
	ExpressionData string // Собственно выражение, которое нужно решить
	CreatedTime    string // Время создания
	LastUpdateTime string // время Последнего апдейта
	Result         string // результат выражения
}

type Agent struct {
	ID       string // ID
	Status   string // Status: Active, Inactive, Lost, Disconnected
	Address  string // Адрес агента
	LastPing string // Время последнего пинга
}

type Task struct {
	ID         string `json:"id"`         // ID
	Expression string `json:"expression"` // Собственно выражение, которое нужно решить
}

func GenerateExpression(expressionData, state string) Expression {
	timeNow := time.Now().Format("2006-01-02 15:04:05.000")
	return Expression{
		Status:         state,
		ExpressionData: expressionData,
		CreatedTime:    timeNow,
		LastUpdateTime: timeNow,
		Result:         fmt.Sprintf("%s = Not Solved", expressionData),
	}
}

func GenerateAgent(agentID, agentAddress string) Agent {
	timeNow := time.Now().Format("2006-01-02 15:04:05.000")
	return Agent{
		ID:       agentID,
		Status:   "Inactive",
		Address:  agentAddress,
		LastPing: timeNow,
	}
}

// Функция для генериции ответов
func Resp(status int, response any) Response {
	return Response{
		Status: status,
		Resp:   response,
	}
}
