package insights

import (
	"net/http"
	"time"
)

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
