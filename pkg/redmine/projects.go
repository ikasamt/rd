package redmine

import (
	"fmt"
	"net/url"
)

type ProjectsResponse struct {
	Projects   []Project `json:"projects"`
	TotalCount int       `json:"total_count"`
	Offset     int       `json:"offset"`
	Limit      int       `json:"limit"`
}

type ProjectResponse struct {
	Project ProjectDetail `json:"project"`
}

type ProjectDetail struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Identifier  string   `json:"identifier"`
	Description string   `json:"description"`
	Status      int      `json:"status"`
	IsPublic    bool     `json:"is_public"`
	Trackers    []Tracker `json:"trackers,omitempty"`
}

func (c *Client) ListProjects() (*ProjectsResponse, error) {
	params := url.Values{}
	params.Set("limit", "100")
	
	var response ProjectsResponse
	if err := c.Get("/projects.json", params, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetProject(id string) (*ProjectDetail, error) {
	path := fmt.Sprintf("/projects/%s.json", id)
	
	var response ProjectResponse
	if err := c.Get(path, nil, &response); err != nil {
		return nil, err
	}

	return &response.Project, nil
}