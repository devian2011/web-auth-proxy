package log

import (
	"fmt"
	"log"
	"time"
)

type StdoutLogger struct {
}

func NewStdoutLogger() *StdoutLogger {
	return &StdoutLogger{}
}

func (sl *StdoutLogger) Write(message Message) {
	msg := sl.buildMessage(message)
	if message.GetLevel() == CriticalLevel || message.GetLevel() == FatalLevel {
		log.Fatalln(msg)
	} else {
		log.Println(msg)
	}
}

func (sl *StdoutLogger) buildMessage(msg Message) string {
	date := time.Now()
	return fmt.Sprintf(
		"[%s] - [%s] Message: %s File: %s:%d Error: %v",
		msg.GetLevel(),
		date.Format("2006-01-02 15:03:04"),
		msg.GetMessage(),
		msg.GetFile(),
		msg.GetLine(),
		msg.GetError())
}

type FileLogger struct {
}

func NewFileLogger() *FileLogger {
	return &FileLogger{}
}

func (fl *FileLogger) Write(message Message) {
	msg := fl.buildMessage(message)
	log.Println(msg)
}

func (fl *FileLogger) buildMessage(msg Message) string {
	date := time.Now()
	return fmt.Sprintf(
		"[%s] - [%s] Message: %s File: %s:%d Error: %v",
		msg.GetLevel(),
		date.Format("2006-01-02 15:03:04"),
		msg.GetMessage(),
		msg.GetFile(),
		msg.GetLine(),
		msg.GetError())
}
