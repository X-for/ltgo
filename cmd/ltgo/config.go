package main

import (
	"fmt"
	"strings"

	"github.com/X-for/ltgo/internal/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long: `Manage ltgo configuration.
If run without arguments, it displays the current configuration.

Available keys: language, site, cookie`,
	Run: func(cmd *cobra.Command, args []string) {
		// 默认行为：显示配置
		showConfig()
	},
}

// configSetCmd 用于修改配置
// 变量名加上 config 前缀，防止和其他文件里的 setCmd 冲突（如果有的话）
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Example: `  ltgo config set language python3
  ltgo config set site com`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		setConfig(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
}

func showConfig() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	fmt.Println("Current Configuration:")
	fmt.Printf("  Language: %s\n", cfg.Language)
	fmt.Printf("  Site:     %s\n", cfg.Site)

	cookiePreview := "Not set"
	if len(cfg.Cookie) > 20 {
		cookiePreview = cfg.Cookie[:20] + "..."
	}
	fmt.Printf("  Cookie:   %s\n", cookiePreview)
}

func setConfig(key, value string) {
	key = strings.ToLower(key)
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	switch key {
	case "language", "lang":
		cfg.Language = value
	case "site":
		if value != "cn" && value != "com" {
			fmt.Println("Error: site must be 'cn' or 'com'")
			return
		}
		cfg.Site = value
	case "cookie":
		cfg.Cookie = value
	default:
		fmt.Printf("Error: unknown configuration key '%s'\n", key)
		return
	}

	if err := cfg.Save(); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		return
	}

	fmt.Printf("✅ Updated %s to '%s'\n", key, value)
}
