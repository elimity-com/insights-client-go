# Elimity Insights Go client

This Go module provides functions for connector interactions with an Elimity Insights server.

## Usage

```go
package main

import (
	"github.com/elimity-com/insights-client-go/v6"
	"time"
)

func main() {
	now := time.Now()
	log := insights.ConnectorLog{
		Level:     "info",
		Message:   "Hello world!",
		Timestamp: now,
	}
	logs := []insights.ConnectorLog{log}
	if err := insights.CreateSourceConnectorLogs("https://example.elimity.com", "42", "my-token", logs); err != nil {
		panic(err)
	}
}
```

## Installation

```sh
$ go get github.com/elimity-com/insights-client-go/v6
```

## Setting up the development environment

### Generating a GitHub access token for accessing the development image

The development image is hosted on the GitHub Container registry, more specifically at
`ghcr.io/elimity-com/insights-client-go`. You can access the registry by
[generating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-token)
with read permission for packages, and then providing it to Docker by running `docker login ghcr.io`.

### Setting up the development environment

This repository contains a `docker-compose.yml` file to configure a development container. You can set up the
environment by running `docker compose up -d`.

## Compatibility

| Client version | Insights version |
|----------------|------------------|
| 1              | 2.7 - 2.11       |
| 2              | 2.12 - 3.0       |
| 3              | 3.1 - 3.5        |
| 4              | ^3.6             |
| 6              | ^3.27            |
