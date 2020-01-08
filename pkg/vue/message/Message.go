package message

type Type string

const (
	SUCCESS = "success"
	WARNING = "warning"
	INFO    = "info"
	ERROR   = "error"
)

type Message struct {
	Message string `json:"message"`
	Type    Type   `json:"type"`
}

func Success(message string) Message {
	return Message{
		Message: message,
		Type:    SUCCESS,
	}
}

func Warning(message string) Message {
	return Message{
		Message: message,
		Type:    WARNING,
	}
}

func Info(message string) Message {
	return Message{
		Message: message,
		Type:    INFO,
	}
}

func Error(message string) Message {
	return Message{
		Message: message,
		Type:    ERROR,
	}
}
