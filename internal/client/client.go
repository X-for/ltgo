package client

import (
	"net/http"
	"time"
	"github.com/X-for/ltgo/internal/config"
)

type Client struct {
	http	 *http.Client
	cfg	 *config.Config
	BaseURL	 string
}

func New(cfg *config.Config) *Client {
	baseURL := "https://leetcode.com"
	if cfg.Site == "cn" {
		baseURL = "https://leetcode.cn"
	}
	return &Client{
		http: &http.Client{
			Timeout: 10 * time.Sec,
		},
		cfg: cfg,
		BaseURL: baseURL,
	}

func Get(url string) ([]byte, error) {
	return http.Get(url)
}

func Post(url string, body []byte) ([]byte, error) {
	return http.Post(url, body)
}

func (c *Client) enhanceRequest(req *http.Request) {
	req.Header.Set("Cookie", c.cfg.Cookie)
	req,Header.Set("User-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
}
