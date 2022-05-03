package log

import (
	json2 "encoding/json"
	"fmt"
	"time"
)

func StrFormatter(msg Message) string {
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

func JsonFormatter(msg Message) string {
	date := time.Now()
	data := struct {
		Date    string `json:"date"`
		Level   string `json:"level"`
		Message string `json:"message"`
		File    string `json:"file"`
		Line    int    `json:"line"`
		Error   string `json:"error"`
	}{
		Date:    date.Format("2006-01-02 15:03:04"),
		Level:   string(msg.GetLevel()),
		Message: msg.GetMessage(),
		File:    msg.GetFile(),
		Line:    msg.GetLine(),
		Error:   msg.GetError().Error(),
	}

	json, _ := json2.Marshal(data)
	return string(json)
}
