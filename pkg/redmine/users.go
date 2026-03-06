package redmine

type UserDetail struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"firstname"`
}

type CurrentUserResponse struct {
	User UserDetail `json:"user"`
}

func (c *Client) GetCurrentUser() (*UserDetail, error) {
	var response CurrentUserResponse
	if err := c.Get("/users/current.json", nil, &response); err != nil {
		return nil, err
	}
	return &response.User, nil
}
