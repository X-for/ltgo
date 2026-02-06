package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/X-for/ltgo/internal/config"
)

type Client struct {
	http     *http.Client
	cfg      *config.Config
	BaseURL  string
	EndPoint string
}

func New(cfg *config.Config) *Client {
	baseURL := "https://leetcode.com"
	endpoint := "https://leetcode.com/graphql"
	if cfg.Site == "cn" {
		baseURL = "https://leetcode.cn"
		endpoint = "https://leetcode.cn/graphql/"
	}
	return &Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		cfg:      cfg,
		BaseURL:  baseURL,
		EndPoint: endpoint,
	}
}

func (c *Client) Get(path string) ([]byte, error) {
	url := c.BaseURL + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.enhanceRequest(req)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) Post(path string, body []byte) ([]byte, error) {
	url := c.BaseURL + path
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.enhanceRequest(req)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) enhanceRequest(req *http.Request) {
	if c.cfg.Cookie != "" {
		req.Header.Set("Cookie", c.cfg.Cookie)
	}
	req.Header.Set("Referer", c.BaseURL+"/problemset/all/")
	req.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
}

// 在 import 里添加 "encoding/json"

// GraphQLPayload 发送给服务器的请求体结构
type GraphQLPayload struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

// GraphQL 发送 GraphQL 请求并解析结果到 target
func (c *Client) GraphQL(query string, variables interface{}, target interface{}) error {
	// 1. 构造请求体
	payload := GraphQLPayload{
		Query:     query,
		Variables: variables,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// 2. 发送 POST 请求 (注意：GraphQL 的 Endpoint 通常是 /graphql)
	// 如果是 .cn 站点，可能是 /graphql/
	respBody, err := c.Post("/graphql/", body)
	if err != nil {
		return err
	}

	fmt.Println("DEBUG:", string(respBody))

	// 3. 解析响应
	if err := json.Unmarshal(respBody, target); err != nil {
		return err
	}

	return nil
}
