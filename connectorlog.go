package insights

import (
	"fmt"
	"net/http"
	"time"
)

// Alertf sends an alert log to the given client's configured server. The log's message is constructed by formatting
// according to the given specifier and arguments, and its timestamp is set to the current time.
func (c Client) Alertf(format string, args ...interface{}) error {
	return c.logf(Alert, format, args...)
}

// CreateConnectorLogs creates the given connector logs at the given client's configured server.
func (c Client) CreateConnectorLogs(logs []ConnectorLog) error {
	pathComponents := []string{"connectorLogs"}
	var models []connectorLog
	for _, log := range logs {
		model := log.model()
		models = append(models, model)
	}
	return c.performRequest(http.MethodPost, pathComponents, models, nil)
}

// Infof sends an info log to the given client's configured server. The log's message is constructed by formatting
// according to the given specifier and arguments, and its timestamp is set to the current time.
func (c Client) Infof(format string, args ...interface{}) error {
	return c.logf(Info, format, args...)
}

func (c Client) logf(level Level, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	timestamp := time.Now()
	log := ConnectorLog{
		Level:     level,
		Message:   message,
		Timestamp: timestamp,
	}
	logs := []ConnectorLog{log}
	return c.CreateConnectorLogs(logs)
}

// ConnectorLog represents a connector log that may be managed at an Elimity Insights server.
type ConnectorLog struct {
	Level     Level
	Message   string
	Timestamp time.Time
}

func (l ConnectorLog) model() connectorLog {
	level := l.Level.model()
	return connectorLog{
		Level:     level,
		Message:   l.Message,
		Timestamp: l.Timestamp,
	}
}

// Level represents a connector log's severity level.
type Level uint

// Valid levels.
const (
	Alert Level = iota
	Info
)

func (l Level) model() string {
	switch l {
	case Alert:
		return "alert"

	case Info:
		return "info"

	default:
		panic("unreachable")
	}
}

type connectorLog struct {
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
