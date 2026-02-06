package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/X-for/ltgo/internal/models"
)

// Generate 生成题目文件到指定目录
// q: 题目详情
// outputDir: 输出目录 (比如 "./questions" 或绝对路径)
func Generate(q *models.QuestionDetail, outputDir string) error {
	// 1. 确保目标目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 2. 构造文件名
	// 格式推荐: 0001_two-sum.go (ID补零方便排序)
	// 如果 ID 是数字，最好补零；如果是面试题(如 "面试题 01.01")，就直接用
	filename := fmt.Sprintf("%s_%s.go", q.QuestionFrontendID, q.TitleSlug)
	fullPath := filepath.Join(outputDir, filename)

	// 3. 提取 Go 代码
	var code string
	for _, s := range q.CodeSnippets {
		if s.LangSlug == "golang" {
			code = s.Code
			break
		}
	}

	if code == "" {
		return fmt.Errorf("no Go code snippet found for question: %s", q.Title)
	}

	// 4. 拼接完整文件内容
	// 这里我们预留位置给后续的注释和 package 声明
	fileContent := fmt.Sprintf(`package main

import "fmt"

/**
 * ID: %s
 * Title: %s
 * Difficulty: %s
 */

%s
`, q.QuestionFrontendID, q.Title, q.Difficulty, code)

	// 5. 检查文件是否已存在 (防止误覆盖)
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("file already exists: %s", fullPath)
	}

	// 6. 写入文件
	fmt.Printf("Generating file: %s\n", fullPath)
	return os.WriteFile(fullPath, []byte(fileContent), 0644)
}
