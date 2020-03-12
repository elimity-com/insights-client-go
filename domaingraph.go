package insights

import "net/http"

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
	pathComponents := []string{"domain-graph", "reload"}
	requestBody := domainGraph.model()
	return c.performRequest(http.MethodPost, pathComponents, requestBody, nil)
}

// DomainGraph represents a graph of domain data that may be managed at an Elimity Insights server.
type DomainGraph struct {
	Entities      []Entity
	Relationships []Relationship
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

	return domainGraphModel{
		Entities:      entityModels,
		Relationships: relationshipModels,
	}
}

// Entity represents an entity that may be managed at an Elimity Insights server.
type Entity struct {
	Active               bool
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
		Active:               e.Active,
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
		FromID:               r.FromEntityID,
		FromType:             r.FromEntityType,
		ToID:                 r.ToEntityID,
		ToType:               r.ToEntityType,
	}
}

type attributeAssignmentModel struct {
	AttributeTypeName string     `json:"attributeTypeName"`
	Value             valueModel `json:"value"`
}

type domainGraphModel struct {
	Entities      []entityModel       `json:"entities"`
	Relationships []relationshipModel `json:"relationships"`
}

type entityModel struct {
	Active               bool                       `json:"active"`
	AttributeAssignments []attributeAssignmentModel `json:"attributeAssignments"`
	ID                   string                     `json:"id"`
	Name                 string                     `json:"name"`
	Type                 string                     `json:"type"`
}

type relationshipModel struct {
	AttributeAssignments []attributeAssignmentModel `json:"attributeAssignments"`
	FromID               string                     `json:"fromId"`
	FromType             string                     `json:"fromType"`
	ToID                 string                     `json:"toId"`
	ToType               string                     `json:"toType"`
}
