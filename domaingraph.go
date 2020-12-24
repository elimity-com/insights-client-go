package insights

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	tim "time"
)

// AttributeAssignment represents the assignment of an attribute of a given type with a compatible value, which may be
// managed at an Elimity Insights server.
//
// A given attribute type and value are compatible if their types are equal. For example, a boolean attribute type can
// only be assigned with a boolean value.
type AttributeAssignment struct {
	AttributeTypeName string
	Value             Value
}

func (a AttributeAssignment) model() attributeAssignmentModel {
	return attributeAssignmentModel{
		AttributeTypeName: a.AttributeTypeName,
		Value:             a.Value.model(),
	}
}

// ReloadDomainGraph reloads the given domain graph at the given client's configured server.
func (c Client) ReloadDomainGraph(domainGraph DomainGraph) error {
	requestBody := new(bytes.Buffer)
	writer := zlib.NewWriter(requestBody)
	encoder := json.NewEncoder(writer)
	model := domainGraph.model()
	if err := encoder.Encode(model); err != nil {
		panic(err)
	}
	if err := writer.Close(); err != nil {
		panic(err)
	}
	return c.performRequest("custom-connector-domain-graphs", "application/octet-stream", requestBody)
}

// DomainGraph represents a graph of domain data that may be managed at an Elimity Insights server.
// The Timestamp field is optional.
type DomainGraph struct {
	Entities      []Entity
	Relationships []Relationship
	Timestamp     *tim.Time
}

func (g DomainGraph) model() domainGraphModel {
	entityModels := []entityModel{}
	for _, entity := range g.Entities {
		entityModel := entity.model()
		entityModels = append(entityModels, entityModel)
	}

	relationshipModels := []relationshipModel{}
	for _, relationship := range g.Relationships {
		relationshipModel := relationship.model()
		relationshipModels = append(relationshipModels, relationshipModel)
	}

	historyTimestamp := parseDomainGraphTimestamp(g.Timestamp)
	return domainGraphModel{
		Entities:         entityModels,
		HistoryTimestamp: historyTimestamp,
		Relationships:    relationshipModels,
	}
}

// Entity represents an entity that may be managed at an Elimity Insights server.
type Entity struct {
	AttributeAssignments []AttributeAssignment
	ID                   string
	Name                 string
	Type                 string
}

func (e Entity) model() entityModel {
	attributeAssignmentModels := []attributeAssignmentModel{}
	for _, attributeAssignment := range e.AttributeAssignments {
		attributeAssignmentModel := attributeAssignment.model()
		attributeAssignmentModels = append(attributeAssignmentModels, attributeAssignmentModel)
	}

	return entityModel{
		AttributeAssignments: attributeAssignmentModels,
		ID:                   e.ID,
		Name:                 e.Name,
		Type:                 e.Type,
	}
}

// Relationship represents a relationship between entities that may be managed at an Elimity Insights server.
type Relationship struct {
	AttributeAssignments []AttributeAssignment
	FromEntityID         string
	FromEntityType       string
	ToEntityID           string
	ToEntityType         string
}

func (r Relationship) model() relationshipModel {
	attributeAssignmentModels := []attributeAssignmentModel{}
	for _, attributeAssignment := range r.AttributeAssignments {
		attributeAssignmentModel := attributeAssignment.model()
		attributeAssignmentModels = append(attributeAssignmentModels, attributeAssignmentModel)
	}

	return relationshipModel{
		AttributeAssignments: attributeAssignmentModels,
		FromEntityID:         r.FromEntityID,
		FromEntityType:       r.FromEntityType,
		ToEntityID:           r.ToEntityID,
		ToEntityType:         r.ToEntityType,
	}
}

type attributeAssignmentModel struct {
	AttributeTypeName string     `json:"attributeTypeName"`
	Value             valueModel `json:"value"`
}

func parseDomainGraphTimestamp(time *tim.Time) dateTime {
	if time == nil {
		return dateTime{}
	}
	return parseDateTime(*time)
}

type domainGraphModel struct {
	Entities         []entityModel       `json:"entities"`
	HistoryTimestamp dateTime            `json:"historyTimestamp"`
	Relationships    []relationshipModel `json:"relationships"`
}

type entityModel struct {
	AttributeAssignments []attributeAssignmentModel `json:"attributeAssignments"`
	ID                   string                     `json:"id"`
	Name                 string                     `json:"name"`
	Type                 string                     `json:"type"`
}

type relationshipModel struct {
	AttributeAssignments []attributeAssignmentModel `json:"attributeAssignments"`
	FromEntityID         string                     `json:"fromEntityId"`
	FromEntityType       string                     `json:"fromEntityType"`
	ToEntityID           string                     `json:"toEntityId"`
	ToEntityType         string                     `json:"toEntityType"`
}
