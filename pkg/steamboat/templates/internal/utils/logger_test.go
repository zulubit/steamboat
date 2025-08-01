package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitLogger(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)
	
	InitLogger()
	
	if Logger == nil {
		t.Fatal("InitLogger did not set global Logger")
	}
	
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		t.Error("logs directory was not created")
	}
	
	Logger.Info("test message")
	if _, err := os.Stat(filepath.Join("logs", "app.log")); os.IsNotExist(err) {
		t.Error("log file was not created after logging")
	}
}

func TestMultipleInitCalls(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)
	
	InitLogger()
	firstLogger := Logger
	
	InitLogger()
	secondLogger := Logger
	
	if firstLogger != secondLogger {
		t.Error("Multiple InitLogger calls should return same instance")
	}
}

func TestLoggerLevels(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)
	
	// Create a fresh logger for this test
	testLogger := newLogger()
	
	testLogger.Info("info message")
	testLogger.Error("error message")
	testLogger.Warn("warn message")
	
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		t.Error("logs directory was not created")
	}
}