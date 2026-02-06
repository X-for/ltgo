package main

import (
	"fmt"
	"log"

	// 这里的路径必须和你 go.mod 里的 module 名字对应
	"github.com/X-for/ltgo/internal/config"
)

func main() {
	// 1. 加载配置 (使用 Load 方法)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err) // 遇到错误直接打印并退出
	}

	// 2. 打印读取到的信息
	fmt.Printf("Current Site: %s\n", cfg.Site)
	fmt.Printf("Current Language: %s\n", cfg.Language)
	fmt.Printf("Old Cookie: %s\n", cfg.Cookie)

	// 3. 修改配置并保存
	cfg.Cookie = "test-cookie-updated-" + cfg.Language
	err = cfg.Save()
	if err != nil {
		log.Fatal("Failed to save config:", err)
	}

	fmt.Println("Config saved successfully!")
}
