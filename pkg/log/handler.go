package log

type Level string

const (
	FatalLevel    Level = "fatal"
	CriticalLevel Level = "critical"
	ErrorLevel    Level = "error"
	WarningLevel  Level = "warning"
	InfoLevel     Level = "info"
	DebugLevel    Level = "debug"
)

type Message interface {
	GetLevel() Level
	GetError() error
	GetLine() int
	GetFile() string
	GetMessage() string
}

type Logger interface {
	Write(msg Message)
}

type Handler struct {
	loggerMap map[Level][]Logger
}

func NewHandler() *Handler {
	return &Handler{
		loggerMap: make(map[Level][]Logger),
	}
}

func (h *Handler) AddLogger(levels []Level, logger Logger) {
	for _, lvl := range levels {
		if _, ok := h.loggerMap[lvl]; !ok {
			h.loggerMap[lvl] = make([]Logger, 0)
		}
		h.loggerMap[lvl] = append(h.loggerMap[lvl], logger)
	}
}

func (h *Handler) Handle(msg Message) {
	if loggers, ok := h.loggerMap[msg.GetLevel()]; ok {
		for _, l := range loggers {
			l.Write(msg)
		}
	}
}
