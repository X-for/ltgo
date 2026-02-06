package models

// Question 列表里的简略信息
type Question struct {
	QuestionFrontendID string `json:"questionFrontendId"`
	Title              string `json:"title"`
	TitleSlug          string `json:"titleSlug"`
	Difficulty         string `json:"difficulty"`
	Status             string `json:"status"`
	IsPaidOnly         bool   `json:"isPaidOnly"`
}

// CodeSnippet 代码模板
type CodeSnippet struct {
	Lang     string `json:"lang"`
	LangSlug string `json:"langSlug"`
	Code     string `json:"code"`
}

// QuestionDetail 题目详细信息
type QuestionDetail struct {
	QuestionFrontendID string        `json:"questionFrontendId"`
	Title              string        `json:"title"`
	TitleSlug          string        `json:"titleSlug"`
	Content            string        `json:"content"`           // 题目描述 (HTML)
	TranslatedContent  string        `json:"translatedContent"` // 中文描述 (CN特有)
	Difficulty         string        `json:"difficulty"`
	CodeSnippets       []CodeSnippet `json:"codeSnippets"` // 各语言代码模板
	SampleTestCase     string        `json:"sampleTestCase"`
}

// QuestionListResponse Data Wrappers 用于 GraphQL 响应解析
type QuestionListResponse struct {
	Data struct {
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
