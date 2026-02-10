package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/X-for/ltgo/internal/models"
)

// getBestContent default Chinese, then English
func getBestContent(q *models.QuestionDetail) string {
	if q.TranslatedContent != "" {
		return q.TranslatedContent
	}
	return q.Content
}

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
	// 选择语言内容
	descHTML := q.Content
	if q.TranslatedContent != "" {
		descHTML = q.TranslatedContent
	}

	// 转换为纯文本
	descText := htmlToText(descHTML)

	// 将文本每一行前面加 "// " 变成注释
	descComment := formatComment(descText)

	// 在code上下加标记
	wrappedCode := fmt.Sprintf("// @lc code=start\n%s\n// @lc code=end", code)

	// 拼接完整文件内容
	fileContent := fmt.Sprintf(`package main

import "fmt"

/**
 * ID: %s
 * Title: %s
 * Difficulty: %s
 *
%s
 */

%s
`, q.QuestionFrontendID, q.Title, q.Difficulty, descComment, wrappedCode)

	// 5. 检查文件是否已存在 (防止误覆盖)
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("file already exists: %s", fullPath)
	}

	// 6. 写入文件
	fmt.Printf("Generating file: %s\n", fullPath)
	return os.WriteFile(fullPath, []byte(fileContent), 0644)
}

func formatComment(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = " * " + line
	}
	return strings.Join(lines, "\n")
}
