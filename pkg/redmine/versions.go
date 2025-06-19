package redmine

import (
	"fmt"
)

type Version struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date,omitempty"`
	CreatedOn   string `json:"created_on"`
	UpdatedOn   string `json:"updated_on"`
}

type VersionsResponse struct {
	Versions []Version `json:"versions"`
}

func (c *Client) ListVersions(projectID string) (*VersionsResponse, error) {
	path := fmt.Sprintf("/projects/%s/versions.json", projectID)
	
	var response VersionsResponse
	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) FindVersionByName(projectID, versionName string) (*Version, error) {
	versions, err := c.ListVersions(projectID)
	if err != nil {
		return nil, err
	}

	for _, version := range versions.Versions {
		if version.Name == versionName {
			return &version, nil
		}
	}

	return nil, fmt.Errorf("version '%s' not found in project '%s'", versionName, projectID)
}