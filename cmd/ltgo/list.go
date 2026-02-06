package main

import (
	"fmt"
	"os"
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

	for _, q := range resp.Data.ProblemsetQuestionList.Questions {
		status := " "
		if q.Status == "ac" {
			status = "✓"
		}

		// 打印每一行，\t 表示换列
		fmt.Fprintf(w, "[%s]\t%s\t%s\t%s\n", status, q.QuestionFrontendID, q.Title, q.Difficulty)
	}

	w.Flush() // 必须 Flush 才能输出
}
