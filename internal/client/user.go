package client

// UserStatus 简单的用户状态模型
type UserStatus struct {
	IsSignedIn bool   `json:"isSignedIn"`
	Username   string `json:"username"`
}

type UserStatusResponse struct {
	Data struct {
		UserStatus UserStatus `json:"userStatus"`
	} `json:"data"`
}

// GetUser 获取当前登录用户信息
func (c *Client) GetUser() (*UserStatus, error) {
	query := `
    query globalData {
        userStatus {
            isSignedIn
            username
        }
    }`

	var resp UserStatusResponse
	if err := c.GraphQL(query, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Data.UserStatus, nil
}
