package models

// Question 题目基础信息
type Question struct {
	QuestionFrontendID string `json:"questionFrontendId"` // 比如 "1"
	Title              string `json:"title"`              // 比如 "Two Sum"
	TitleSlug          string `json:"titleSlug"`          // 比如 "two-sum"
	Difficulty         string `json:"difficulty"`         // 比如 "Easy"
	IsPaidOnly         bool   `json:"isPaidOnly"`
}

// QuestionListResponse 题目列表的响应结构
type QuestionListResponse struct {
	Data struct {
		ProblemsetQuestionList struct {
			Total     int        `json:"total"`
			Questions []Question `json:"questions"`
		} `json:"problemsetQuestionList"`
	} `json:"data"`
}
