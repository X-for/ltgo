package generator

import (
	"regexp"
	"strings"
)

// htmlToText 将 HTML 转换为适合放在 Go 注释里的纯文本
func htmlToText(html string) string {
	// 1. 替换常见块级标签为换行，保证段落清晰
	// 把 </p>, </div>, <br> 换成换行符
	reBlock := regexp.MustCompile(`(</p>|</div>|<br/?>)`)
	text := reBlock.ReplaceAllString(html, "\n")

	// 2. 去掉所有剩余的 HTML 标签 (<...>)
	reTag := regexp.MustCompile(`<[^>]*>`)
	text = reTag.ReplaceAllString(text, "")

	// 3. 处理常见的 HTML 实体 (Entities)
	// 标准库 html.UnescapeString 可以做这个，但这里手动处理几个最常见的也行
	replacer := strings.NewReplacer(
		"&nbsp;", " ",
		"&lt;", "<",
		"&gt;", ">",
		"&amp;", "&",
		"&quot;", "\"",
		"&#39;", "'",
	)
	text = replacer.Replace(text)

	// 4. 清理多余的空行 (连续换行变成一个)
	reNewlines := regexp.MustCompile(`\n{3,}`)
	text = reNewlines.ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}
