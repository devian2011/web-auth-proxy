package log

type LoggableMessage struct {
	level   Level
	err     error
	line    int
	file    string
	message string
}

func (em LoggableMessage) GetLevel() Level {
	return em.level
}

func (em LoggableMessage) GetError() error {
	return em.err
}

func (em LoggableMessage) GetLine() int {
	return em.line
}

func (em LoggableMessage) GetFile() string {
	return em.file
}

func (em LoggableMessage) GetMessage() string {
	return em.message
}

func NewLoggableMessage(message string, line int, file string, err error, level Level) LoggableMessage {
	return LoggableMessage{
		level:   level,
		err:     err,
		line:    line,
		file:    file,
		message: message,
	}
}

func NewErrorMessage(message string, line int, file string, err error) LoggableMessage {
	return NewLoggableMessage(message, line, file, err, ErrorLevel)
}

func NewCriticalMessage(message string, line int, file string, err error) LoggableMessage {
	return NewLoggableMessage(message, line, file, err, CriticalLevel)
}

func NewFatalMessage(message string, line int, file string, err error) LoggableMessage {
	return NewLoggableMessage(message, line, file, err, FatalLevel)
}

func NewWarningMessage(message string, line int, file string, err error) LoggableMessage {
	return NewLoggableMessage(message, line, file, err, WarningLevel)
}

func NewInfoMessage(message string, line int, file string, err error) LoggableMessage {
	return NewLoggableMessage(message, line, file, err, InfoLevel)
}

func NewDebugMessage(message string, line int, file string, err error) LoggableMessage {
	return NewLoggableMessage(message, line, file, err, DebugLevel)
}
