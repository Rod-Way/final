package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"task/app/api"
	"task/app/validator"
	db "task/database"

	"go.mongodb.org/mongo-driver/mongo"
)

// Добавление вычисления арифметического выражения.
func AddExp(ctx context.Context, log *slog.Logger, col *mongo.Collection, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")              // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
		}

		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
			return
		}

		// Чтение POST запроса
		body, err := io.ReadAll(r.Body)
		if err != nil && err != io.EOF {
			http.Error(w, fmt.Sprintf("body read error: %v", err), http.StatusBadRequest)
			return
		}
		exp := string(body)

		// Проверка выражения
		var isOK bool = validator.IsValid(exp)

		var state string
		if isOK {
			state = "Active"
		} else {
			state = "Bad"
		}

		expression := api.GenerateExpression(exp, state)

		// Добавляем информацию в базу данных
		if err = db.Add(ctx, col, expression); err != nil {
			log.Error("failed add expression to bd")
			http.Error(w, fmt.Sprintf("failed add expression to bd: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		res := api.Resp(http.StatusOK, exp)
		fmt.Fprint(w, res)
	}
}

/*===========================================================================================================================================*/
/*===========================================================================================================================================*/

// Получение списка выражений со статусами.
func GetAllExp(ctx context.Context, log *slog.Logger, col *mongo.Collection, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")             // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
		}

		res, err := db.GetAll(ctx, col)
		if err != nil {
			http.Error(w, fmt.Sprintf("MongoDB error: %v", err), http.StatusInternalServerError)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON encoding error: %v", err), http.StatusInternalServerError)
		}
	}
}

// Получение значения выражения по его идентификатору.
func GetExpByID(ctx context.Context, log *slog.Logger, col *mongo.Collection, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")             // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
		}

		// Получаем значение id
		id := r.URL.Query().Get("id")

		if id == "" {
			http.Error(w, "Bad id value", http.StatusBadRequest)
		}

		res, err := db.GetByID(ctx, col, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("MongoDB error: %v", err), http.StatusInternalServerError)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON encoding error: %v", err), http.StatusInternalServerError)
		}
	}
}

/*===========================================================================================================================================*/
/*===========================================================================================================================================*/

// Получение списка доступных операций со временем их выполения.
func GetOperations(ctx context.Context, log *slog.Logger, col *mongo.Collection, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")             // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
		}

		// Получаем значение операции
		operation := r.URL.Query().Get("operation")

		// Массив длительностей операций
		var res []api.OperationDuration

		// Получаем длительность операций
		switch operation {
		case "all":
			if operations, err := db.GetOperations(ctx, col); err != nil {
				http.Error(w, fmt.Sprintf("Error: %d", err), http.StatusBadRequest)
			} else {
				res = operations
			}
			break
		case "plus", "minus", "multiply", "divide", "agent":
			res = append(res, api.OperationDuration{OperationName: operation, OperationDuration: 10})
			break
		default:
			http.Error(w, "Invalid 'operation' parameter", http.StatusBadRequest)
			return
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON encoding error: %v", err), http.StatusInternalServerError)
		}
	}
}

func SetOperations(ctx context.Context, log *slog.Logger, col *mongo.Collection, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")              // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
		}

		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
			return
		}

		// Декодируем
		var operations []api.OperationDuration
		if err := json.NewDecoder(r.Body).Decode(&operations); err != nil {
			http.Error(w, fmt.Sprintf("JSON decoding error: %v", err), http.StatusBadRequest)
			return
		}

		// Можно было бы захендлить больше ошибок, но -TODO: вставить отмазку-
		if len(operations) != 5 {
			http.Error(w, fmt.Sprint("Operation lenght error!"), http.StatusBadRequest)
			return
		}
		var isOK = true
		for _, oper := range operations {
			if oper.OperationDuration <= 0 {
				http.Error(w, fmt.Sprint("Operation Duration error: Duration is equals to zero or negative nummber"), http.StatusBadRequest)
				isOK = false
				break
			}
		}

		// OK?
		if isOK {
			if err := db.SetOperations(ctx, col, operations); err != nil {
				http.Error(w, fmt.Sprintf("Error: %d", err), http.StatusBadRequest)
			}

			w.WriteHeader(http.StatusOK)
		} else if !isOK {
			w.WriteHeader(http.StatusBadRequest)
		}
		exp := "I love Pizza"
		res := api.Resp(http.StatusOK, exp)
		fmt.Fprint(w, res)
	}

}

/*===========================================================================================================================================*/
/*===========================================================================================================================================*/

func RegisterAgent(ctx context.Context, log *slog.Logger, col *mongo.Collection, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")              // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
		}

		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Получение значения поля ID
		id := r.Form.Get("ID")

		// Получение значения поля Address
		address := r.Form.Get("Address")

		// Проверка выражения
		// var isOK bool = validator.IsValid(exp)
		isOK := true
		// OK?
		if isOK {

			agent := api.GenerateAgent(id, address)

			// Добавляем информацию в базу данных
			if err = db.AddAgent(ctx, col, agent); err != nil {
				log.Error("failed add expression to bd")
				http.Error(w, fmt.Sprintf("failed add expression to bd: %v", err), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		} else if !isOK {
			w.WriteHeader(http.StatusBadRequest)
		}

		res := api.Resp(http.StatusOK, id)
		fmt.Fprint(w, res)
	}
}

func GetAgents(ctx context.Context, log *slog.Logger, col *mongo.Collection, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")             // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
		}

		agents, err := db.GetAllAgents(ctx, col)
		if err != nil {
			http.Error(w, fmt.Sprintf("MongoDB error: %v", err), http.StatusInternalServerError)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(agents)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON encoding error: %v", err), http.StatusInternalServerError)
		}

	}
}

func UpdateAgent(ctx context.Context, log *slog.Logger, col *mongo.Collection, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")              // Разрешить запросы от всех доменов
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Разрешить GET и OPTIONS запросы
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Разрешить заголовок Content-Type

		// Проверяем метод запроса
		if r.Method == "OPTIONS" {
			// Отправляем пустой ответ для предварительного запроса OPTIONS
			w.WriteHeader(http.StatusOK)
		}

		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
			return
		}

		var agent api.Agent
		if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
			http.Error(w, fmt.Sprintf("JSON decoding error: %v", err), http.StatusBadRequest)
			return
		}

		err := db.UpdateAgent(ctx, col, agent)
		if err != nil {
			http.Error(w, fmt.Sprintf("MongoDB error: %v", err), http.StatusInternalServerError)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode("Good")
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON encoding error: %v", err), http.StatusInternalServerError)
		}
	}
}

/*===========================================================================================================================================*/
/*===========================================================================================================================================*/

// Получение задачи для выполения.
func GetTask(ctx context.Context, log *slog.Logger, col *mongo.Collection, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
			return
		}
	}
}

// Приём результата обработки данных.
func SubmitResult(ctx context.Context, log *slog.Logger, col *mongo.Collection, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
			return
		}
	}
}
