package insights

import "net/http"

// AttributeType represents an attribute type that may be managed at an Elimity Insights server.
type AttributeType struct {
	Category    string
	Description string
	Name        string
	Type        Type
}

func (t AttributeType) model() attributeTypeModel {
	typeModel := t.Type.model()
	return attributeTypeModel{
		Category:    t.Category,
		Description: t.Description,
		Name:        t.Name,
		Type:        typeModel,
	}
}

// CreateAttributeType creates the given attribute type at the given client's configured server.
func (c Client) CreateAttributeType(attributeType AttributeType) error {
	pathComponents := []string{"attributeTypes"}
	requestBody := attributeType.model()
	return c.performRequest(http.MethodPost, pathComponents, requestBody, nil)
}

type attributeTypeModel struct {
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	Type        typeModel `json:"type"`
}
