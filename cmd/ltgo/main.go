package main

import (
    "fmt"
    "log"

    "github.com/X-for/ltgo/internal/client"
    "github.com/X-for/ltgo/internal/config"
)

func main() {
    cfg, _ := config.Load()
    // 强制设置 site 为 cn 确保环境
    cfg.Site = "cn"
    lc := client.New(cfg)

    // 这个 Query 专门获取题目详情，不需要登录也能拿
    query := `
    query questionData($titleSlug: String!) {
        question(titleSlug: $titleSlug) {
            questionId
            title
difficulty
        }
    }`

    vars := map[string]interface{}{
        "titleSlug": "two-sum",
    }

    // 使用 map[string]interface{} 来接收任意结构的响应，避免 struct 定义错误导致解析失败
    var resp map[string]interface{}

    fmt.Println("Sending request to LeetCode CN...")
    err := lc.GraphQL(query, vars, &resp)
    if err != nil {
        log.Fatal(err)
    }

    // 此时 client.go 里的调试打印会把原始 JSON 吐出来
}