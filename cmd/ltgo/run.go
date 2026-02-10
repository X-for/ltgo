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

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Run code on LeetCode",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startRun(args[0])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func startRun(filePath string) {
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

	// 6. 提交运行
	fmt.Printf("🚀 Sending code (%s) to LeetCode...\n", lang)
	interpretID, err := c.RunCode(q, code, lang)
	if err != nil {
		fmt.Printf("Failed to submit run: %v\n", err)
		return
	}

	// 7. 轮询结果
	fmt.Print("Waiting for result...")
	res, err := c.CheckResult(interpretID)
	if err != nil {
		fmt.Printf("\nError checking result: %v\n", err)
		return
	}
	fmt.Println("\n")

	// 8. 漂亮地打印结果
	// 编译错误
	if res.CompileError != "" || res.FullCompileError != "" {
		fmt.Println("❌ Compile Error:")
		if res.FullCompileError != "" {
			fmt.Println(res.FullCompileError)
		} else {
			fmt.Println(res.CompileError)
		}
		return
	}

	// 运行时错误
	if res.RuntimeError != "" {
		fmt.Println("❌ Runtime Error:")
		fmt.Println(res.RuntimeError)
		return
	}

	// [修改 1] 使用 TotalTestcases 来控制循环，防止出现空的 Case 2
	count := res.TotalTestcases
	if count == 0 {
		count = len(res.CodeAnswer)
	}

	// 打印总结
	if (res.StatusMsg == "Accepted" || res.StatusMsg == "Finished") && res.CorrectAnswer {
		fmt.Println("✅ Accepted\n")
	} else if res.StatusMsg == "Compile Error" {
		// ... (其实前面已经拦截了编译错误)
	} else {
		// 其他情况统统算 Wrong Answer (只要代码跑完了但 CorrectAnswer 是 false)
		fmt.Println("❌ Wrong Answer\n")
		// 如果想看原始状态，可以保留: fmt.Printf("(Status: %s)\n", res.StatusMsg)
	}

	// 详细打印每个 Case
	for i := 0; i < count; i++ {
		input := ""
		// [修改 2] 尝试获取 Input，如果 API 没返回，就显示 SampleTestCase
		if i < len(res.InputFormatted) {
			input = res.InputFormatted[i]
		} else if i == 0 && q.SampleTestCase != "" {
			// 只有第一个 case 我们能确信是 SampleTestCase
			// 简单的格式化一下，把换行符换成空格，避免太长
			input = strings.ReplaceAll(q.SampleTestCase, "\n", " ")
		}

		output := ""
		if i < len(res.CodeAnswer) {
			output = res.CodeAnswer[i]
		}

		expected := ""
		if i < len(res.ExpectedOutput) {
			expected = res.ExpectedOutput[i]
		}

		stdOut := ""
		if i < len(res.StdOutput) && res.StdOutput[i] != "" {
			stdOut = res.StdOutput[i]
		}

		fmt.Printf("Case %d:\n", i+1)
		if input != "" {
			fmt.Printf("  Input:    %s\n", input)
		}
		fmt.Printf("  Output:   %s\n", output)
		if expected != "" {
			fmt.Printf("  Expected: %s\n", expected)
		}
		if stdOut != "" {
			fmt.Printf("  Stdout:   %s\n", stdOut)
		}
		fmt.Println("  ------------------------")
	}

}
