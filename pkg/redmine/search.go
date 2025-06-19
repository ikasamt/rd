package redmine

import (
	"net/url"
	"strconv"
)

// SearchResult represents a single search result
type SearchResult struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Datetime    string `json:"datetime"`
}

// SearchResponse represents the response from /search.json
type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	TotalCount int            `json:"total_count"`
	Offset     int            `json:"offset"`
	Limit      int            `json:"limit"`
}

// SearchOptions represents search parameters
type SearchOptions struct {
	Query       string
	Offset      int
	Limit       int
	Scope       string // all, my_projects, subprojects
	AllWords    bool
	TitlesOnly  bool
	Issues      bool
	News        bool
	Documents   bool
	Changesets  bool
	WikiPages   bool
	Messages    bool
	Projects    bool
}

// Search performs a search using Redmine's /search API
func (c *Client) Search(opts *SearchOptions) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("q", opts.Query)
	
	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Scope != "" {
		params.Set("scope", opts.Scope)
	}
	
	// Set resource type filters
	if opts.Issues {
		params.Set("issues", "1")
	}
	if opts.News {
		params.Set("news", "1")
	}
	if opts.Documents {
		params.Set("documents", "1")
	}
	if opts.Changesets {
		params.Set("changesets", "1")
	}
	if opts.WikiPages {
		params.Set("wiki_pages", "1")
	}
	if opts.Messages {
		params.Set("messages", "1")
	}
	if opts.Projects {
		params.Set("projects", "1")
	}
	
	if opts.AllWords {
		params.Set("all_words", "1")
	}
	if opts.TitlesOnly {
		params.Set("titles_only", "1")
	}
	
	var response SearchResponse
	if err := c.Get("/search.json", params, &response); err != nil {
		return nil, err
	}
	
	return &response, nil
}