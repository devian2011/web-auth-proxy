package log

import (
	"log"
)


type Formatter func(message Message) string

type StdoutLogger struct {
	formatter Formatter
}

type Color string

const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed          = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorBlue         = "\u001b[34m"
	ColorReset        = "\u001b[0m"
)

func (sl *StdoutLogger) colorize(color Color, message string) {
	log.Println(string(color), message, string(ColorReset))
}

func NewStdoutLogger(formatter Formatter) *StdoutLogger {
	return &StdoutLogger{
		formatter: formatter,
	}
}

func (sl *StdoutLogger) Write(message Message) {
	msg := sl.formatter(message)
	if message.GetLevel() == CriticalLevel || message.GetLevel() == FatalLevel {
		sl.colorize(ColorRed, msg)
	} else {
		log.Println(msg)
	}
}

type FileLoggerConfig struct {
	Path string `json:"path"`
	Level []Level `json:"level"`
}

type FileLogger struct {
	formatter Formatter
}

func NewFileLogger(formatter Formatter) *FileLogger {
	return &FileLogger{
		formatter: formatter,
	}
}

func (fl *FileLogger) Write(message Message) {
	log.Println(fl.formatter(message))
}

