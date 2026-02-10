package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/X-for/ltgo/internal/client"
	"github.com/X-for/ltgo/internal/config"
	"github.com/X-for/ltgo/internal/generator"
	"github.com/X-for/ltgo/internal/models"
	"github.com/spf13/cobra"
)

func isNumeric(s string) bool {
	match, _ := regexp.MatchString(`^\d+$`, s)
	return match
}

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

	// 1. 尝试调用服务端搜索
	matches, err := c.SearchQuestions(arg)
	if err != nil {
		fmt.Printf("Search failed: %v\n", err)
		return
	}

	// 2. 如果没找到，或者用户输入的是纯数字 ID，我们可以尝试本地兜底
	// (因为服务端搜索搜 ID 有时候不准)
	// 但为了简化，我们先只用服务端结果。

	if len(matches) == 0 {
		fmt.Println("❌ No questions found.")
		return
	}

	var targetQ models.Question
	foundExact := false

	// 3. 智能匹配：先看有没有 ID 或 Slug 完全一样的
	for _, q := range matches {
		if q.QuestionFrontendID == arg || q.TitleSlug == arg {
			targetQ = q
			foundExact = true
			break
		}
	}

	if foundExact {
		// 找到了精确匹配，直接开始
		fmt.Printf("🎯 Found exact match: [%s] %s\n", targetQ.QuestionFrontendID, targetQ.Title)
	} else if len(matches) == 1 {
		// 只有一个模糊匹配结果，也就它了
		targetQ = matches[0]
		fmt.Printf("🎯 Found: [%s] %s\n", targetQ.QuestionFrontendID, targetQ.Title)
	} else {
		// 多个结果，列出来让用户选
		fmt.Println("Multiple questions found:")
		for _, q := range matches {
			fmt.Printf(" - [%s] %s\n", q.QuestionFrontendID, q.Title)
		}
		fmt.Println("\n⚠️  Please use the exact ID or Slug to generate.")
		return
	}

	// 4. 获取详情并生成
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
