package redmine

import (
	"fmt"
	"net/url"
	"strconv"
)

type IssueFilter struct {
	ProjectID  string
	StatusID   string
	AssignedTo string
	Limit      int
	Offset     int
}

func (c *Client) ListIssues(filter *IssueFilter) (*IssuesResponse, error) {
	params := url.Values{}
	
	if filter != nil {
		if filter.ProjectID != "" {
			params.Set("project_id", filter.ProjectID)
		}
		if filter.StatusID != "" {
			params.Set("status_id", filter.StatusID)
		}
		if filter.AssignedTo != "" {
			params.Set("assigned_to_id", filter.AssignedTo)
		}
		if filter.Limit > 0 {
			params.Set("limit", strconv.Itoa(filter.Limit))
		} else {
			params.Set("limit", "25")
		}
		if filter.Offset > 0 {
			params.Set("offset", strconv.Itoa(filter.Offset))
		}
	} else {
		params.Set("limit", "25")
	}

	var response IssuesResponse
	if err := c.Get("/issues.json", params, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetIssue(id int, includeJournals bool) (*Issue, error) {
	path := fmt.Sprintf("/issues/%d.json", id)
	
	params := url.Values{}
	if includeJournals {
		params.Set("include", "journals")
	}

	var response IssueResponse
	if err := c.Get(path, params, &response); err != nil {
		return nil, err
	}

	return &response.Issue, nil
}