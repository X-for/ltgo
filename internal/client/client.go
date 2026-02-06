package client

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/X-for/ltgo/internal/config"
)

type Client struct {
	http    *http.Client
	cfg     *config.Config
	BaseURL string
}

func New(cfg *config.Config) *Client {
	baseURL := "https://leetcode.com"
	if cfg.Site == "cn" {
		baseURL = "https://leetcode.cn"
	}
	return &Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		cfg:     cfg,
		BaseURL: baseURL,
	}
}

func (c *Client) Get(url string) ([]byte, error) {
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

func (c *Client) Post(url string, body []byte) ([]byte, error) {
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
	req.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
}
