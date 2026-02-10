package generator

import (
	"fmt"
	"os"
	"path/filepath"

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
// outputDir: 输出目录
// site: "cn" 或 "com"
// lang: 目标语言 slug (e.g. "golang", "python3")
func Generate(q *models.QuestionDetail, outputDir string, site string, lang string) error {
	// 1. 获取语言配置
	langConf := GetLangConfig(lang)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 2. 构造文件名 (使用正确的后缀)
	filename := fmt.Sprintf("%s_%s.%s", q.QuestionFrontendID, q.TitleSlug, langConf.Extension)
	fullPath := filepath.Join(outputDir, filename)

	// 3. 提取对应语言的代码 Snippet
	var code string
	for _, s := range q.CodeSnippets {
		if s.LangSlug == lang { // <--- 动态匹配
			code = s.Code
			break
		}
	}

	if code == "" {
		// 如果没找到指定语言，尝试回退到 Go 或者报错
		// 这里直接报错比较好，提示用户
		return fmt.Errorf("no code snippet found for language: %s", lang)
	}

	// 4. 准备注释内容
	descHTML := q.Content
	if q.TranslatedContent != "" {
		descHTML = q.TranslatedContent
	}
	descText := htmlToText(descHTML)

	// 格式化注释 (根据语言风格)
	var descComment string
	var metaBlock string

	appName := "leetcode.com"
	if site == "cn" {
		appName = "leetcode.cn"
	}

	// Code Marker
	markerStart := fmt.Sprintf("%s @lc code=start", langConf.Comment)
	markerEnd := fmt.Sprintf("%s @lc code=end", langConf.Comment)
	wrappedCode := fmt.Sprintf("%s\n%s\n%s", markerStart, code, markerEnd)

	// 构造注释块
	if langConf.Comment == "//" {
		// C-style (Go, Java, C++, JS, TS, Rust)
		// 使用 /** ... */ 包裹
		metaBlock = fmt.Sprintf(`/*
 * @lc app=%s id=%s lang=%s
 * @lc slug=%s
 * @lc type=question
 */`, appName, q.QuestionFrontendID, lang, q.TitleSlug)

		commentBody := FormatComment(descText, langConf)
		descComment = fmt.Sprintf(`/**
 * ID: %s
 * Title: %s
 * Difficulty: %s
 *
%s
 */`, q.QuestionFrontendID, q.Title, q.Difficulty, commentBody)

	} else {
		// Script-style (Python, Ruby, Shell)
		// 使用 # 逐行注释
		metaBlock = fmt.Sprintf(`%s @lc app=%s id=%s lang=%s
%s @lc slug=%s
%s @lc type=question`, langConf.Comment, appName, q.QuestionFrontendID, lang, langConf.Comment, q.TitleSlug, langConf.Comment)

		descComment = fmt.Sprintf(`%s ID: %s
%s Title: %s
%s Difficulty: %s
%s
%s`, langConf.Comment, q.QuestionFrontendID, langConf.Comment, q.Title, langConf.Comment, q.Difficulty, langConf.Comment, FormatComment(descText, langConf))
	}

	// 5. 拼接完整文件内容
	var fileContent string

	if lang == "golang" {
		// Go 特殊处理: 需要 package main 和 import
		fileContent = fmt.Sprintf("package main\n\nimport \"fmt\"\n\n%s\n\n%s\n\n%s\n", metaBlock, descComment, wrappedCode)
	} else {
		// 其他语言直接拼接
		fileContent = fmt.Sprintf("%s\n\n%s\n\n%s\n", metaBlock, descComment, wrappedCode)
	}

	// 6. 检查文件是否已存在
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("file already exists: %s", fullPath)
	}

	// 7. 写入文件
	fmt.Printf("Generating file: %s\n", fullPath)
	return os.WriteFile(fullPath, []byte(fileContent), 0644)
}
