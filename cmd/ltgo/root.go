package main

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd 代表基础命令，也就是直接运行 `ltgo` 时执行的
var rootCmd = &cobra.Command{
	Use:   "ltgo",
	Short: "A CLI tool for LeetCode",
	Long:  `ltgo is a fast and simple CLI tool for LeetCode in Go.`,
	// 这里通常留空，或者是显示帮助信息
}

// Execute 是 main.go 调用的入口
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// 在这里可以定义全局 flag，比如 --site cn
	// rootCmd.PersistentFlags().StringVar(...)
}
