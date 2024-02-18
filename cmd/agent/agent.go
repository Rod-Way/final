package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"task/cmd/agent/parser"
	"task/config"
)

var (
	maxAgents, maxGoroutines int
	wg                       sync.WaitGroup
)

func main() {
	logger := setupLogger()
	cfg := config.MustLoad(*logger)
	str := cfg.HTTPServer.Host
	logger.Info(str)
	maxGoroutines, _ = strconv.Atoi(cfg.GoroutinesNum)
	maxAgents, _ = strconv.Atoi(cfg.AgentsNum)

	for i := 0; i < maxAgents; i++ {
		wg.Add(1)
		go agent()
	}

	wg.Wait()
}

func agent() {
	defer wg.Done()
	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exp := "2+2+2" // Пример выражения, можно подставить любое другое
			result, _ := parser.Calculate(exp)
			fmt.Printf("Expression: %s\nResult: %d\n\n", exp, result)
		}()
	}
}

func setupLogger() *slog.Logger {
	logg := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logg
}
