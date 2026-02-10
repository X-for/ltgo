package generator

import "strings"

type LangConfig struct {
	Slug      string
	Extension string
	Comment   string // 单行注释符
}

// SupportedLangs 支持的语言配置
// 可以在这里添加更多语言支持
var SupportedLangs = map[string]LangConfig{
	"golang":     {Slug: "golang", Extension: "go", Comment: "//"},
	"python3":    {Slug: "python3", Extension: "py", Comment: "#"},
	"java":       {Slug: "java", Extension: "java", Comment: "//"},
	"cpp":        {Slug: "cpp", Extension: "cpp", Comment: "//"},
	"c":          {Slug: "c", Extension: "c", Comment: "//"},
	"javascript": {Slug: "javascript", Extension: "js", Comment: "//"},
	"typescript": {Slug: "typescript", Extension: "ts", Comment: "//"},
	"rust":       {Slug: "rust", Extension: "rs", Comment: "//"},
}

func GetLangConfig(lang string) LangConfig {
	if cfg, ok := SupportedLangs[lang]; ok {
		return cfg
	}
	// 默认回退到 Go
	return SupportedLangs["golang"]
}

// FormatComment 根据语言生成多行注释块
func FormatComment(text string, langConf LangConfig) string {
	lines := strings.Split(text, "\n")
	var sb strings.Builder

	// 统一采用每行加单行注释符的风格，最稳健
	// C系: // Line 1
	// Py:  # Line 1
	for _, line := range lines {
		// 为了美观，C系语言通常习惯用 /** ... */ 块注释
		// 但为了简化逻辑，统一用单行注释符也是完全合法的

		// 如果是 C 系语言 (//)，我们可以保留原来的 * 风格，只要外层包裹好
		if langConf.Comment == "//" {
			sb.WriteString(" * " + line + "\n")
		} else {
			sb.WriteString(langConf.Comment + " " + line + "\n")
		}
	}
	return strings.TrimSpace(sb.String())
}
