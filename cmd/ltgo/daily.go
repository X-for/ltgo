package main

import (
	"fmt"
	"os"

	"github.com/X-for/ltgo/internal/client"
	"github.com/X-for/ltgo/internal/config"
	"github.com/X-for/ltgo/internal/generator"
	"github.com/spf13/cobra"
)

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Get today's daily question",
	Long:  `Fetch and generate the daily question from LeetCode.`,
	Run: func(cmd *cobra.Command, args []string) {
		runDaily()
	},
}

func init() {
	rootCmd.AddCommand(dailyCmd)
}

func runDaily() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Please run 'ltgo init' first.")
		return
	}
	c := client.New(cfg)

	fmt.Println("Fetching daily question...")
	q, err := c.GetDailyQuestion()
	if err != nil {
		fmt.Printf("Failed to get daily question: %v\n", err)
		return
	}

	fmt.Printf("📅 Today's Question: [%s] %s (%s)\n", q.QuestionFrontendID, q.Title, q.Difficulty)

	// 复用生成逻辑
	fmt.Printf("Fetching details for '%s'...\n", q.TitleSlug)
	detail, err := c.GetQuestionDetail(q.TitleSlug)
	if err != nil {
		fmt.Printf("Failed to get details: %v\n", err)
		return
	}

	cwd, _ := os.Getwd()
	outputDir := fmt.Sprintf("%s/questions", cwd)

	if err := generator.Generate(detail, outputDir, cfg.Site, cfg.Language); err != nil {
		fmt.Printf("Failed to generate: %v\n", err)
		return
	}

	fmt.Println("Done! Happy Coding! 🚀")
}
