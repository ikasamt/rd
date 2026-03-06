package redmine

import "fmt"

type CustomFieldDefinition struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CustomFieldsResponse struct {
	CustomFields []CustomFieldDefinition `json:"custom_fields"`
}

func (c *Client) ListCustomFields() ([]CustomFieldDefinition, error) {
	var response CustomFieldsResponse
	if err := c.Get("/custom_fields.json", nil, &response); err != nil {
		return nil, fmt.Errorf("failed to list custom fields (admin permission required): %w", err)
	}
	return response.CustomFields, nil
}

// FindCustomFieldByName はカスタムフィールド名からIDを解決する
func (c *Client) FindCustomFieldByName(name string) (*CustomFieldDefinition, error) {
	fields, err := c.ListCustomFields()
	if err != nil {
		return nil, err
	}
	for _, f := range fields {
		if f.Name == name {
			return &f, nil
		}
	}
	return nil, fmt.Errorf("custom field '%s' not found", name)
}
