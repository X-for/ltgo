package models

// Question 基础题目信息 (适配 V2)
type Question struct {
	QuestionFrontendID string `json:"questionFrontendId"`
	Title              string `json:"title"`
	TranslatedTitle    string `json:"translatedTitle"` // 新增: 中文标题
	TitleSlug          string `json:"titleSlug"`
	Difficulty         string `json:"difficulty"` // "EASY", "MEDIUM", "HARD"
	Status             string `json:"status"`     // "TO_DO", "AC", "TRIED" (可能为null)
	PaidOnly           bool   `json:"paidOnly"`   // 注意: JSON 里是 paidOnly
	IsPaidOnly         bool   `json:"isPaidOnly"` // 兼容旧版
}

// CodeSnippet 代码模板
type CodeSnippet struct {
	Lang     string `json:"lang"`
	LangSlug string `json:"langSlug"`
	Code     string `json:"code"`
}

// QuestionDetail 题目详细信息
type QuestionDetail struct {
	QuestionID         string        `json:"questionId"`
	QuestionFrontendID string        `json:"questionFrontendId"`
	Title              string        `json:"title"`
	TitleSlug          string        `json:"titleSlug"`
	Content            string        `json:"content"`           // 题目描述 (HTML)
	TranslatedContent  string        `json:"translatedContent"` // 中文描述 (CN特有)
	Difficulty         string        `json:"difficulty"`
	CodeSnippets       []CodeSnippet `json:"codeSnippets"` // 各语言代码模板
	SampleTestCase     string        `json:"sampleTestCase"`
}

// QuestionListResponse 题目列表的响应结构
type QuestionListResponse struct {
	Data struct {
		// 兼容 V2
		ProblemsetQuestionListV2 struct {
			Total     int        `json:"total"` // 如果 V2 不返回这个，可能就是 0
			Questions []Question `json:"questions"`
		} `json:"problemsetQuestionListV2"`

		// 保留旧版兼容 (可选)
		ProblemsetQuestionList struct {
			Total     int        `json:"total"`
			Questions []Question `json:"questions"`
		} `json:"problemsetQuestionList"`
	} `json:"data"`
}

type QuestionDetailResponse struct {
	Data struct {
		Question QuestionDetail `json:"question"`
	} `json:"data"`
}
