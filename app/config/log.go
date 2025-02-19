package config

import (
	"log"
	"log/slog"
	"os"
)

func InitLogger(logFilePath string, level slog.Level) (*slog.Logger, func()) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть файл для логов: %v", err)
	}
	handler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	cleanup := func() {
		file.Close()
	}
	return logger, cleanup
}
