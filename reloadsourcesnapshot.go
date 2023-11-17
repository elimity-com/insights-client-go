package insights

import (
	"compress/zlib"
	"encoding/json"
	"io"
)

// ReloadSourceSnapshot sends the given entities and relationships to the referenced Elimity Insights server.
func ReloadSourceSnapshot(
	entities []Entity, insightsURL, sourceID, sourceToken string, relationships []Relationship,
	skipSSLVerification bool,
) error {
	reader, writer := io.Pipe()
	snapshot := snapshot{
		Entities:      entities,
		Relationships: relationships,
	}
	go writeSnapshot(snapshot, writer)
	err := request(
		"application/octet-stream", insightsURL, sourceID, sourceToken, "%s/api/sources/%s/snapshots", reader,
		skipSSLVerification,
	)
	_ = reader.Close()
	return err
}

func writeSnapshot(snapshot snapshot, pipeWriter *io.PipeWriter) {
	zlibWriter := zlib.NewWriter(pipeWriter)
	encoder := json.NewEncoder(zlibWriter)
	_ = encoder.Encode(snapshot)
	_ = zlibWriter.Close()
	_ = pipeWriter.Close()
}

// AttributeAssignment represents the assignment of a value for an attribute type to an entity or a relationship.
type AttributeAssignment struct {
	AttributeTypeID string
	Value           Value
}

// Entity represents a domain graph entity together with its attribute assignments.
type Entity struct {
	AttributeAssignments []AttributeAssignment
	ID                   string
	Name                 string
	Type                 string
}

// Relationship represents a domain graph relationship together with its attribute assignments.
type Relationship struct {
	AttributeAssignments []AttributeAssignment
	FromEntityID         string
	FromEntityType       string
	ToEntityID           string
	ToEntityType         string
}

// Value represents a typed value of an attribute assignment.
type Value struct {
	Type  string
	Value any
}

type snapshot struct {
	Entities      []Entity
	Relationships []Relationship
}
