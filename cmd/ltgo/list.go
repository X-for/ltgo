package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter" // 用这个对齐输出，超好用

	"github.com/X-for/ltgo/internal/client"
	"github.com/X-for/ltgo/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List questions",
	Run: func(cmd *cobra.Command, args []string) {
		runList()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Please run 'ltgo init' first.")
		return
	}

	c := client.New(cfg)

	// 2. 获取题目 (我们先写死获取前 50 题，后面可以加 flag 控制)
	fmt.Println("Fetching questions...")
	resp, err := c.GetQuestions(50, 0)
	if err != nil {
		fmt.Printf("Failed to fetch questions: %v\n", err)
		return
	}

	// 3. 格式化输出
	// tabwriter 可以自动对齐列
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// 表头
	fmt.Fprintln(w, "Status\tID\tTitle\tDifficulty")
	fmt.Fprintln(w, "------\t--\t-----\t----------")

	// 适配 V2 和旧版数据
	questions := resp.Data.ProblemsetQuestionListV2.Questions
	if len(questions) == 0 {
		questions = resp.Data.ProblemsetQuestionList.Questions
	}

	// 遍历 questions 打印
	for _, q := range questions {
		status := " "
		// 状态码转换 (V2 返回的是 TO_DO / AC)
		if q.Status == "ac" || q.Status == "AC" {
			status = "✓"
		}

		// 难度首字母大写转换 (EASY -> Easy)
		diff := q.Difficulty
		if len(diff) > 1 {
			diff = diff[0:1] + strings.ToLower(diff[1:])
		}

		// 优先显示中文标题 (如果有)
		title := q.Title
		if q.TranslatedTitle != "" {
			title = fmt.Sprintf("%s (%s)", q.TranslatedTitle, q.Title)
		}

		fmt.Fprintf(w, "[%s]\t%s\t%s\t%s\n", status, q.QuestionFrontendID, title, diff)
	}

	w.Flush()
}
