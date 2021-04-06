# Elimity Insights Go client

This Go module provides a client for connector interactions with an Elimity Insights server.

## Usage

```go
package main

import (
	"github.com/elimity-com/insights-client-go/v5"
	"time"
)

func main() {
	sourceID := 4
	client, err := insights.NewClient("https://local.elimity.com:8081/api", "token", sourceID)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	log := insights.ConnectorLog{
		Level:     insights.Info,
		Message:   "Hello world!",
		Timestamp: now,
	}
	logs := []insights.ConnectorLog{log}
	if err := client.CreateConnectorLogs(logs); err != nil {
		panic(err)
	}
}
```

## Installation

```sh
$ go get github.com/elimity-com/insights-client-go/v5
```

## Compatibility

| Client version | Insights version |
| -------------- | ---------------- |
| 1              | 2.7 - 2.11       |
| 2              | 2.12 - 3.0       |
| 3              | 3.1 - 3.5        |
| 4              | ^3.6             |
