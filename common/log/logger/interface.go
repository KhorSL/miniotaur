package logger

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

type noop struct{}

func NewNoop() Logger {
	return &noop{}
}

func (n *noop) Debug(msg string, fields map[string]interface{}) {}
func (n *noop) Info(msg string, fields map[string]interface{})  {}
func (n *noop) Warn(msg string, fields map[string]interface{})  {}
func (n *noop) Error(msg string, fields map[string]interface{}) {}
func (n *noop) Fatal(msg string, fields map[string]interface{}) {}
