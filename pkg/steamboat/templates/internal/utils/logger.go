package utils

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *slog.Logger
	once   sync.Once
)

func InitLogger() {
	once.Do(func() {
		Logger = newLogger()
	})
}

func newLogger() *slog.Logger {
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic("failed to create logs directory: " + err.Error())
	}

	logFile := &lumberjack.Logger{
		Filename:   filepath.Join("logs", "app.log"),
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     0,
		Compress:   false,
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler)
}