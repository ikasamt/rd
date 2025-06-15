package redmine

import "time"

type Issue struct {
	ID             int                    `json:"id"`
	Project        Project                `json:"project"`
	Tracker        Tracker                `json:"tracker"`
	Status         Status                 `json:"status"`
	Priority       Priority               `json:"priority"`
	Author         User                   `json:"author"`
	AssignedTo     *User                  `json:"assigned_to,omitempty"`
	Subject        string                 `json:"subject"`
	Description    string                 `json:"description"`
	StartDate      *string                `json:"start_date,omitempty"`
	DueDate        *string                `json:"due_date,omitempty"`
	DoneRatio      int                    `json:"done_ratio"`
	EstimatedHours *float64               `json:"estimated_hours,omitempty"`
	CustomFields   []CustomField          `json:"custom_fields,omitempty"`
	CreatedOn      time.Time              `json:"created_on"`
	UpdatedOn      time.Time              `json:"updated_on"`
	Journals       []Journal              `json:"journals,omitempty"`
}

type IssuesResponse struct {
	Issues     []Issue `json:"issues"`
	TotalCount int     `json:"total_count"`
	Offset     int     `json:"offset"`
	Limit      int     `json:"limit"`
}

type IssueResponse struct {
	Issue Issue `json:"issue"`
}

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tracker struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Priority struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CustomField struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type Journal struct {
	ID        int       `json:"id"`
	User      User      `json:"user"`
	Notes     string    `json:"notes"`
	CreatedOn time.Time `json:"created_on"`
	Details   []Detail  `json:"details,omitempty"`
}

type Detail struct {
	Property string `json:"property"`
	Name     string `json:"name"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}