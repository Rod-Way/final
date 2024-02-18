package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	config "task/config"
	db "task/database"
	handlers "task/server/handlers"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Обрабатываем конфиг и создаем логгер
	logger := setupLogger()
	cfg := config.MustLoad(*logger)
	str := cfg.HTTPServer.Host
	logger.Info(str)
	logger.Debug("logger debg mode enabled")
	ctx := context.Background()
	client := &http.Client{}

	mongoCollections := db.NewClient(context.TODO(), logger)

	mux := http.NewServeMux()

	mux.Handle("/expression/add", middlewareSetup(handlers.AddExp(ctx, logger, mongoCollections["Expressions"], client)))

	mux.Handle("/expressions/get-all", middlewareSetup(handlers.GetAllExp(ctx, logger, mongoCollections["Expressions"], client)))
	mux.Handle("/expressions/get-id", middlewareSetup(handlers.GetExpByID(ctx, logger, mongoCollections["Expressions"], client)))

	mux.Handle("/operations/get", middlewareSetup(handlers.GetOperations(ctx, logger, mongoCollections["Operations"], client)))
	mux.Handle("/operations/add", middlewareSetup(handlers.SetOperations(ctx, logger, mongoCollections["Operations"], client)))

	mux.Handle("/agents/register", middlewareSetup(handlers.RegisterAgent(ctx, logger, mongoCollections["Agents"], client)))
	mux.Handle("/agents/get", middlewareSetup(handlers.GetAgents(ctx, logger, mongoCollections["Agents"], client)))
	mux.Handle("/agents/update", middlewareSetup(handlers.UpdateAgent(ctx, logger, mongoCollections["Agents"], client)))

	mux.Handle("/task/get", middlewareSetup(handlers.GetTask(ctx, logger, mongoCollections["Tasks"], client)))
	mux.Handle("/result/get", middlewareSetup(handlers.SubmitResult(ctx, logger, mongoCollections["Expressions"], client)))

	// Запускаем сервер
	fmt.Printf("Server started on port %s", cfg.HTTPServer.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPServer.Port, mux))
}

// Добавляем middleware к хэндлеру
func middlewareSetup(handler http.HandlerFunc) http.Handler {
	return (middleware.Recoverer(
		middleware.Logger(
			handler,
		)))
}

// Создаем сетап для логгера
func setupLogger() *slog.Logger {
	logg := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logg
}
