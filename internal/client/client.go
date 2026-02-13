package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// enhanceRequest 增强请求头，伪装成浏览器以绕过 Cloudflare
func (c *Client) enhanceRequest(req *http.Request) {
	if c.cfg.Cookie != "" {
		req.Header.Set("Cookie", c.cfg.Cookie)
	}

	// 基础 Header
	req.Header.Set("Origin", c.BaseURL)
	req.Header.Set("Referer", c.BaseURL+"/")

	// 浏览器指纹 Header (模仿 Chrome 122)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")

	// Sec-CH-UA 系列 (重要)
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="122", "Not(A:Brand";v="24", "Google Chrome";v="122"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	// 尝试提取并设置 x-csrftoken (CN 站必需)
	if csrftoken := extractCSRF(c.cfg.Cookie); csrftoken != "" {
		req.Header.Set("x-csrftoken", csrftoken)
	}
}

// extractCSRF 从 Cookie 字符串中提取 csrftoken
func extractCSRF(cookie string) string {
	parts := strings.Split(cookie, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "csrftoken=") {
			return strings.TrimPrefix(part, "csrftoken=")
		}
	}
	return ""
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
		// [新增] 只有出错时才打印 Body 前 200 个字符，防止刷屏
		preview := string(respBody)
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		return fmt.Errorf("json parse error: %v, response preview: %s", err, preview)
	}

	return nil
}
