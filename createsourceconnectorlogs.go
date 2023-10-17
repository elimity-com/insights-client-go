package insights

import (
	"bytes"
	"encoding/json"
	"time"
)

// CreateSourceConnectorLogs sends the given logs to the referenced Elimity Insights server.
func CreateSourceConnectorLogs(insightsURL, sourceID, sourceToken string, logs []ConnectorLog) error {
	bys, _ := json.Marshal(logs)
	reader := bytes.NewReader(bys)
	return request("application/json", insightsURL, sourceID, sourceToken, "%s/api/sources/%s/connector-logs", reader)
}

// ConnectorLog represents a timestamped log providing additional information during the import process.
type ConnectorLog struct {
	Level     string
	Message   string
	Timestamp time.Time
}
