package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	InfoLevel = iota
	WarningLevel
	ErrorLevel
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
)

type Logger struct {
	Level       int
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

var logger *Logger

func init() {
	logger = &Logger{
		Level:       InfoLevel,
		infoLogger:  log.New(os.Stdout, "", log.LstdFlags),
		warnLogger:  log.New(os.Stdout, "", log.LstdFlags),
		errorLogger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func SetLevel(level int) {
	logger.Level = level
}

func Info(module, method, message string) {
	if logger.Level <= InfoLevel {
		logger.infoLogger.Println(colorize(colorGreen, formatMessage("INFO", module, method, message)))
	}
}

func Warning(module, method, message string) {
	if logger.Level <= WarningLevel {
		logger.warnLogger.Println(colorize(colorYellow, formatMessage("WARN", module, method, message)))
	}
}

func Error(module, method string, err error, message ...string) {
	if logger.Level <= ErrorLevel {
		fullMessage := err.Error()
		if len(message) > 0 && message[0] != "" {
			fullMessage = fmt.Sprintf("%s | %s", message[0], err.Error())
		}
		logger.errorLogger.Println(colorize(colorRed, formatMessage("ERROR", module, method, fullMessage)))
	}
}

func formatMessage(level, module, method, message string) string {
	return fmt.Sprintf("%s: [Module: %s] [Method: %s] %s", level, module, method, message)
}

func colorize(color, message string) string {
	return fmt.Sprintf("%s%s%s", color, message, colorReset)
}
