package generator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ReadSolution 读取文件并提取 @lc code 标记之间的内容
func ReadSolution(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	text := string(content)
	// 稍微放宽一点匹配条件，有时候用户可能会手动改
	startTag := "@lc code=start"
	endTag := "@lc code=end"

	startIdx := strings.Index(text, startTag)
	endIdx := strings.LastIndex(text, endTag) // 用 LastIndex 防止有多个匹配混淆

	if startIdx == -1 || endIdx == -1 {
		// 尝试旧版格式 (// @lc code=start)
		// 其实上面的 Index 应该也能匹配到，因为包含关系
		return "", fmt.Errorf("code markers (@lc code=start/end) not found in %s", filePath)
	}

	// 调整索引，跳过标签本身
	// startTag 所在行的下一行开始
	// 这里简单处理：找到 startTag 后换行符的位置
	lineEndAfterStart := strings.Index(text[startIdx:], "\n")
	if lineEndAfterStart == -1 {
		return "", fmt.Errorf("invalid code block format")
	}
	realStart := startIdx + lineEndAfterStart + 1

	if realStart >= endIdx {
		return "", fmt.Errorf("empty code block")
	}

	code := text[realStart:endIdx]
	return strings.TrimSpace(code), nil
}

// ParseSlugFromMeta 从文件内容的元数据中提取 Slug
func ParseSlugFromMeta(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	text := string(content)

	// 匹配 @lc slug=xxx
	// 兼容多种写法，只要包含这个 pattern 即可
	re := regexp.MustCompile(`@lc\s+slug=([a-zA-Z0-9-]+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("slug metadata not found")
}

func ParseLangFromMeta(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	text := string(content)
	// 匹配 @lc ... lang=xxx
	re := regexp.MustCompile(`@lc\s+.*lang=([a-zA-Z0-9]+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("lang metadata not found")
}
