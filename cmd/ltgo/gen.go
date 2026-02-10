package main

import (
	"fmt"
	"os"

	"github.com/X-for/ltgo/internal/client"
	"github.com/X-for/ltgo/internal/config"
	"github.com/X-for/ltgo/internal/generator"
	"github.com/X-for/ltgo/internal/models"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen [slug]",
	Short: "Generate a question file",
	Long:  `Generate a Go file for a specific question. Example: ltgo gen two-sum`,
	Args:  cobra.ExactArgs(1), // 必须接受 1 个参数
	Run: func(cmd *cobra.Command, args []string) {
		runGen(args[0])
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}

func runGen(arg string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Please run 'ltgo init' first.")
		return
	}
	c := client.New(cfg)

	fmt.Printf("Searching for '%s'...\n", arg)

	// 1. 搜索题目
	matches, err := c.SearchQuestionsByKeyword(arg)
	if err != nil {
		fmt.Printf("Search failed: %v\n", err)
		return
	}

	if len(matches) == 0 {
		fmt.Println("❌ No questions found.")
		return
	}

	var targetQ models.Question

	// 2. 判断结果
	if len(matches) == 1 {
		// 只有一个命中，直接选中
		targetQ = matches[0]
		fmt.Printf("🎯 Found: [%s] %s\n", targetQ.QuestionFrontendID, targetQ.Title)
	} else {
		// 多个命中，列出来让用户选
		fmt.Println("Multiple questions found:")
		for _, q := range matches {
			fmt.Printf(" - [%s] %s\n", q.QuestionFrontendID, q.Title)
		}
		fmt.Println("\n⚠️  Please refine your search or use the exact ID.")
		return
	}

	// 3. 获取详情并生成
	fmt.Printf("Fetching details for '%s'...\n", targetQ.TitleSlug)
	detail, err := c.GetQuestionDetail(targetQ.TitleSlug)
	if err != nil {
		fmt.Printf("Failed to get details: %v\n", err)
		return
	}

	cwd, _ := os.Getwd()
	outputDir := fmt.Sprintf("%s/questions", cwd)

	if err := generator.Generate(detail, outputDir); err != nil {
		fmt.Printf("Failed to generate: %v\n", err)
		return
	}

	fmt.Println("Done! Happy Coding! 🚀")
}
