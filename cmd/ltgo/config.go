package main

import (
	"fmt"
	"strings"

	"github.com/X-for/ltgo/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Manage ltgo configuration (language, site, etc.)",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}
		fmt.Printf("Language: %s\n", cfg.Language)
		fmt.Printf("Site:     %s\n", cfg.Site)
		// Cookie 比较长且敏感，可以只显示一部分或者不显示
		cookiePreview := "Not set"
		if len(cfg.Cookie) > 20 {
			cookiePreview = cfg.Cookie[:20] + "..."
		}
		fmt.Printf("Cookie:   %s\n", cookiePreview)
	},
}

var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Example: `  ltgo config set language python3
  ltgo config set site cn
  ltgo config set site com`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := strings.ToLower(args[0])
		value := args[1]

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
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(listCmd)
	configCmd.AddCommand(setCmd)
}
