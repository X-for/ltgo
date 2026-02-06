package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/X-for/ltgo/internal/client"
	"github.com/X-for/ltgo/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config and login",
	Long:  "Setup your LeetCode account by inputting cookie.",
	Run: func(cmd *cobra.Command, args []string) {
		runInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit() {
	reader := bufio.NewReader(os.Stdin)

	// 1. 询问站点
	fmt.Print("Choose site (cn/com) [default: cn]: ")
	site, _ := reader.ReadString('\n')
	site = strings.TrimSpace(site)
	if site == "" {
		site = "cn"
	}

	// 2. 询问 Cookie
	fmt.Println("Please paste your LeetCode Cookie (from browser developer tools):")
	fmt.Println("(Include LEETCODE_SESSION and csrftoken)")
	fmt.Print("> ")
	cookie, _ := reader.ReadString('\n')
	cookie = strings.TrimSpace(cookie)

	fmt.Println("\nVerifying your cookie...")

	// 3. 验证 Cookie
	tempCfg := &config.Config{
		Site:     site,
		Cookie:   cookie,
		Language: "golang", // 默认
	}

	c := client.New(tempCfg)
	user, err := c.GetUser()
	if err != nil {
		fmt.Printf("❌ Connection failed: %v\n", err)
		return
	}

	if !user.IsSignedIn {
		fmt.Println("❌ Cookie is invalid or expired (Not Signed In).")
		fmt.Println("Please try again with a fresh cookie.")
		return
	}

	fmt.Printf("✅ Welcome, %s!\n", user.Username)

	// 4. 保存
	if err := tempCfg.Save(); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return
	}

	fmt.Println("🎉 Configuration saved successfully!")
}
