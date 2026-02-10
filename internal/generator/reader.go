package generator

import (
	"fmt"
	"os"
	"strings"
)

// ReadSolution 读取文件并提取 @lc code 标记之间的内容
func ReadSolution(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	text := string(content)
	startTag := "// @lc code=start"
	endTag := "// @lc code=end"

	startIdx := strings.Index(text, startTag)
	endIdx := strings.Index(text, endTag)

	if startIdx == -1 || endIdx == -1 {
		return "", fmt.Errorf("code markers not found in %s", filePath)
	}

	// 提取中间的内容
	code := text[startIdx+len(startTag) : endIdx]
	return strings.TrimSpace(code), nil
}
