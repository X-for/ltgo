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

var (
	difficulty string
	status     string
	tag        string
	id         string
)

var genCmd = &cobra.Command{
	Use:   "gen [keyword]", // 改一下 usage 提示
	Short: "Generate a question file",
	Long: `Generate a Go file for a specific question.
Example: 
  ltgo gen two-sum
  ltgo gen sum --difficulty=Hard
  ltgo gen --tag=dp --status=todo (列出没做的 DP 题)`,
	Args: cobra.MaximumNArgs(1), // 允许不传 keyword，只要有 flag
	Run: func(cmd *cobra.Command, args []string) {
		keyword := ""
		if len(args) > 0 {
			keyword = args[0]
		}
		runGen(keyword)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVarP(&difficulty, "difficulty", "d", "", "Difficulty (Easy, Medium, Hard)")
	genCmd.Flags().StringVarP(&status, "status", "s", "", "Status (todo, solved, attempted)")
	genCmd.Flags().StringVarP(&tag, "tag", "t", "", "Topic Tag (e.g. array, dp)")
	genCmd.Flags().StringVarP(&id, "id", "i", "", "Search by exact Frontend ID")
}

func isNumeric(s string) bool {
	match, _ := regexp.MatchString(`^\d+$`, s)
	return match
}

func runGen(keyword string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Please run 'ltgo init' first.")
		return
	}
	c := client.New(cfg)

	fmt.Printf("Searching for '%s'...\n", keyword)

	// [修改 1] 改用服务端搜索 SearchQuestions (而不是本地 SearchQuestionsByKeyword)
	// 新代码：先构造 Options 结构体
	opts := client.SearchOptions{
		Keyword:    keyword,    // 这里的 keyword 就是原来的 arg
		Difficulty: difficulty, // 需要在 gen.go 里定义这些 flag 变量
		Status:     status,
		Tag:        tag,
		FrontendID: id,
	}
	matches, err := c.SearchQuestions(opts)

	if err != nil {
		fmt.Printf("Search failed: %v\n", err)
		return
	}

	if len(matches) == 0 {
		fmt.Println("❌ No questions found.")
		return
	}

	var targetQ models.Question
	foundExact := false

	// [修改 2] 增加智能匹配逻辑
	// 如果找到了完全匹配的 ID 或 Slug，就不用让用户选了
	for _, q := range matches {
		if q.QuestionFrontendID == keyword || q.TitleSlug == keyword {
			targetQ = q
			foundExact = true
			break
		}
	}

	if foundExact {
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
		fmt.Println("\n⚠️  Please refine your search or use the exact ID.")
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

	if err := generator.Generate(detail, outputDir, cfg.Site, cfg.Language); err != nil {
		fmt.Printf("Failed to generate: %v\n", err)
		return
	}

	fmt.Println("Done! Happy Coding! 🚀")
}
