package insights

import "net/http"

// RelationshipAttributeType represents a relationship attribute type that may be managed at an Elimity Insights server.
type RelationshipAttributeType struct {
	FromCategory string
	Description  string
	Name         string
	ToCategory   string
	Type         Type
}

func (t RelationshipAttributeType) model() relationshipAttributeTypeModel {
	typeModel := t.Type.model()
	return relationshipAttributeTypeModel{
		ChildType:   t.ToCategory,
		Description: t.Description,
		Name:        t.Name,
		ParentType:  t.FromCategory,
		Type:        typeModel,
	}
}

// CreateRelationshipAttributeType creates the given relationship attribute type at the given client's configured
// server.
func (c Client) CreateRelationshipAttributeType(relationshipAttributeType RelationshipAttributeType) error {
	pathComponents := []string{"relationshipAttributeTypes"}
	requestBody := relationshipAttributeType.model()
	return c.performRequest(http.MethodPost, pathComponents, requestBody, nil)
}

type relationshipAttributeTypeModel struct {
	ChildType   string    `json:"childType"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	ParentType  string    `json:"parentType"`
	Type        typeModel `json:"type"`
}
