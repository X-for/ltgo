package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/X-for/ltgo/internal/client"
	"github.com/X-for/ltgo/internal/config"
	"github.com/X-for/ltgo/internal/generator"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit [file]",
	Short: "Submit code to LeetCode",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startSubmit(args[0])
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}

func startSubmit(filePath string) {

	// 1. 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("File not found: %s\n", filePath)
		return
	}

	// 2. 尝试解析 Slug
	var slug string
	// 先尝试从文件元数据里读
	metaSlug, err := generator.ParseSlugFromMeta(filePath)
	if err == nil && metaSlug != "" {
		slug = metaSlug
		// fmt.Printf("Found slug from metadata: %s\n", slug)
	} else {
		// 读不到(旧文件)则回退到文件名解析
		filename := filepath.Base(filePath)
		parts := strings.Split(filename, "_")
		if len(parts) >= 2 {
			slugWithExt := parts[1]
			slug = strings.TrimSuffix(slugWithExt, ".go")
		} else {
			fmt.Println("Could not parse slug from metadata or filename (expected ID_slug.go).")
			return
		}
	}

	// 获取编码语言
	lang, err := generator.ParseLangFromMeta(filePath)
	if err != nil || lang == "" {
		// 如果没找到元数据，尝试根据后缀推断 (兼容旧文件或手写文件)
		ext := strings.TrimPrefix(filepath.Ext(filePath), ".")
		// 简单的反向查找
		for k, v := range generator.SupportedLangs {
			if v.Extension == ext {
				lang = k
				break
			}
		}
		if lang == "" {
			lang = "golang" // 最后的保底
		}
	}

	// 3. 读取代码
	code, err := generator.ReadSolution(filePath)
	if err != nil {
		fmt.Printf("Failed to read solution: %v\n", err)
		return
	}

	// 4. 初始化 Client
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Please run 'ltgo init' first.")
		return
	}
	c := client.New(cfg)

	// 5. 获取题目详情 (为了拿 Test Case 和 ID)
	fmt.Printf("Fetching question info for '%s'...\n", slug)
	q, err := c.GetQuestionDetail(slug)
	if err != nil {
		fmt.Printf("Failed to get question info: %v\n", err)
		return
	}

	// 6. 提交代码
	fmt.Printf("🚀 Submitting to LeetCode...\n")
	subID, err := c.SubmitCode(q, code, lang)
	if err != nil {
		fmt.Printf("Failed to submit: %v\n", err)
		return
	}
	fmt.Printf("Submission ID: %d\n", subID)

	// 7. 轮询结果
	fmt.Print("Waiting for result...")
	res, err := c.CheckSubmission(subID)
	if err != nil {
		fmt.Printf("\nError checking result: %v\n", err)
		return
	}
	fmt.Println("\n")

	// 8. 打印结果
	if res.CompileError != "" {
		fmt.Println("❌ Compile Error:")
		fmt.Println(res.FullCompileError)
		return
	}

	if res.RuntimeError != "" {
		fmt.Println("❌ Runtime Error:")
		fmt.Println(res.RuntimeError)
		return
	}

	if res.StatusMsg == "Accepted" {
		fmt.Println("✅ Accepted!")
		fmt.Printf("Runtime: %s (Beats %.2f%%)\n", res.StatusRuntime, res.RuntimePercentile)
		fmt.Printf("Memory:  %s (Beats %.2f%%)\n", res.StatusMemory, res.MemoryPercentile)
	} else if res.StatusMsg == "Wrong Answer" {
		fmt.Println("❌ Wrong Answer")
		fmt.Printf("Passed:   %d/%d cases\n", res.TotalCorrect, res.TotalTestcases)
		fmt.Printf("Input:    %s\n", res.InputFormatted)
		fmt.Printf("Output:   %s\n", res.CodeOutput)
		fmt.Printf("Expected: %s\n", res.ExpectedOutput)
		if res.StdOutput != "" {
			fmt.Printf("Stdout:   %s\n", res.StdOutput)
		}
	} else {
		fmt.Printf("Status: %s\n", res.StatusMsg)
		fmt.Printf("Passed: %d/%d cases\n", res.TotalCorrect, res.TotalTestcases)
		// 比如 Time Limit Exceeded 时，InputFormatted 也会有值
		if res.InputFormatted != "" {
			fmt.Printf("Last Input: %s\n", res.InputFormatted)
		}
	}
}
